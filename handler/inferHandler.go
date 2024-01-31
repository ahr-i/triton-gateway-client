package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/ahr-i/triton-gateway-client/models"
	"github.com/ahr-i/triton-gateway-client/setting"
	"github.com/ahr-i/triton-gateway-client/src/errController"
	"github.com/ahr-i/triton-gateway-client/src/networkController"
	"github.com/gorilla/mux"
)

/* Response Struct */
type TritonResponse struct {
	Key      string `json:"key"`
	Response string `json:"response"`
}

type TritonOutput struct {
	Outputs []struct {
		Data []float32 `json:"data"`
	} `json:"outputs"`
}

type ResponseData struct {
	Image string `json:"image"`
}

type SchedulerResponse struct {
	Key string `json:"key"`
}

/* Request Struct */
type RequestData struct {
	Prompt string `json:"prompt"`
}

/* Inference Handler: Triton Server에 Inference Request 및 Image 전달 */
func (h *Handler) inferHandler(w http.ResponseWriter, r *http.Request) {
	_, fp, _, _ := runtime.Caller(1)

	// Request Decode
	var request RequestData
	err_ := json.NewDecoder(r.Body).Decode(&request)
	errController.ErrorCheck(err_, "REQUEST JSON DECODE ERROR", fp)
	defer r.Body.Close()

	vars := mux.Vars(r)
	model := vars["name"]

	// log.Println(request.Prompt)
	if request.Prompt == "" || request.Prompt == " " {
		rend.JSON(w, http.StatusBadRequest, nil)

		return
	}

	// Model, Version Check And Setting
	modelMap := models.GetModelList()
	version, err := modelMap[model]
	if !err {
		rend.JSON(w, http.StatusNotFound, nil)

		return
	}

	// Triton Inference Request
	img, err := requestSchedulerAndGetResult(request, model, version)
	if err {
		log.Println("Scheduler ERROR")

		rend.JSON(w, http.StatusBadRequest, nil)
		return
	}

	err_ = saveImageToLocal(img)
	if err_ != nil {
		panic(err)
	}

	imgBase64 := encodeImageInBase64(img)

	rend.JSON(w, http.StatusOK, ResponseData{Image: imgBase64})
}

/* Send an inference request to the Scheduler-node and return the request */
func requestSchedulerAndGetResult(request RequestData, model string, version string) (*image.RGBA, bool) {
	_, fp, _, _ := runtime.Caller(1)

	rand.Seed(time.Now().UnixNano())

	seed := rand.Intn(10001)
	listenPort, err_ := networkController.GetAvailablePort()
	if err_ != nil {
		log.Println("PORT ERROR")

		return nil, true
	}

	urlQuery := "?model=" + model + "&version=" + version + "&address=" + networkController.GetLocalIP() + ":" + strconv.Itoa(listenPort)
	url := "http://" + setting.SchedulerUrl + "/request" + urlQuery
	requestData := map[string]interface{}{
		"inputs": []map[string]interface{}{
			{
				"name":     "PROMPT",
				"datatype": "BYTES",
				"shape":    []int{1},
				"data":     []string{request.Prompt},
			},
			{
				"name":     "SAMPLES",
				"datatype": "INT32",
				"shape":    []int{1},
				"data":     []int{1},
			},
			{
				"name":     "STEPS",
				"datatype": "INT32",
				"shape":    []int{1},
				"data":     []int{45},
			},
			{
				"name":     "GUIDANCE_SCALE",
				"datatype": "FP32",
				"shape":    []int{1},
				"data":     []float32{7.5},
			},
			{
				"name":     "SEED",
				"datatype": "INT64",
				"shape":    []int{1},
				"data":     []int{seed},
			},
		},
		"outputs": []map[string]string{
			{
				"name": "IMAGES",
			},
		},
	}

	requestJSON, err_ := json.Marshal(requestData)
	errController.ErrorCheck(err_, "JSON MARSHAL ERROR", fp)

	req, err_ := http.NewRequest("POST", url, bytes.NewBuffer(requestJSON))
	errController.ErrorCheck(err_, "HTTP REQUEST ERROR", fp)
	req.Header.Set("Content-Type", "application/json")

	// Scheduler Server Response
	client := &http.Client{}
	resp, err_ := client.Do(req)
	errController.ErrorCheck(err_, "HTTP RESPONSE ERROR", fp)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err_ := ioutil.ReadAll(resp.Body)
		errController.ErrorCheck(err_, "RESPONSE READ ERROR", fp)

		var schedulerResponse SchedulerResponse
		if err := json.Unmarshal(body, &schedulerResponse); err != nil {
			log.Fatalf("RESPONSE JSON PARSE ERROR: %v", err)
		}
		log.Println(schedulerResponse.Key)

		img := listenAndGetImage(listenPort, schedulerResponse.Key)

		return img, false
	} else {
		return nil, true
	}
}

func listenAndGetImage(port int, authCode string) *image.RGBA {
	log.Println("Open Port: ", port)
	receiverAddr, err := net.ResolveTCPAddr("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		panic(err)
	}

	receiverConn, err := net.ListenTCP("tcp", receiverAddr)
	if err != nil {
		panic(err)
	}
	defer receiverConn.Close()

	buffer := make([]byte, 17485760)
	for {
		conn, err := receiverConn.AcceptTCP()
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		n, err_ := conn.Read(buffer)
		if err_ != nil {
			panic(err)
		}

		var tritonResponse TritonResponse
		if err := json.Unmarshal(buffer[:n], &tritonResponse); err != nil {
			log.Println("ERR", err)
		}
		log.Println(tritonResponse.Key)

		if tritonResponse.Key == authCode {
			var tritonOutput TritonOutput

			err := json.Unmarshal([]byte(tritonResponse.Response), &tritonOutput)
			if err != nil {
				panic(err)
			}

			img, err_ := converUint8ToPng(tritonOutput.Outputs[0].Data)
			if err_ {
				panic(err)
			}

			return img
		} else {
			log.Println("ERROR")

			conn.Close()
		}
	}
}

func converUint8ToPng(imgData []float32) (*image.RGBA, bool) {
	// Uint8 Array To Image
	if len(imgData) <= 0 {
		return nil, true
	}

	// Image의 크기 가정
	width, height := 512, 512
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// ImgData에서 픽셀 값 추출 및 Image 생성
	for i := 0; i < len(imgData); i += 3 {
		x := (i / 3) % width
		y := (i / 3) / width
		r := uint8(imgData[i] * 255)
		g := uint8(imgData[i+1] * 255)
		b := uint8(imgData[i+2] * 255)
		img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
	}

	return img, false
}

/* Save image to local */
func saveImageToLocal(img *image.RGBA) error {
	currentTime := time.Now().Format("20060102-150405.999")
	fileName := "result-" + currentTime + ".png"
	file, err := os.Create("./result/" + fileName)
	if err != nil {
		log.Fatalf("이미지 파일 생성 실패: %v", err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return err
	}

	return nil
}

/* Encode the image in base64 */
func encodeImageInBase64(img *image.RGBA) string {
	var buffer bytes.Buffer

	if err := png.Encode(&buffer, img); err != nil {
		log.Println("BASE ENCODE FAIL")

		os.Exit(1)
	}
	imgBase64 := base64.StdEncoding.EncodeToString(buffer.Bytes())

	return imgBase64
}
