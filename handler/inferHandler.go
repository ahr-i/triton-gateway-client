package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/ahr-i/triton-gateway-client/setting"
	"github.com/ahr-i/triton-gateway-client/src/networkController"
	"github.com/gorilla/mux"
)

/* Response Struct */
type TritonResponse struct {
	Token    string `json:"token"`
	Response string `json:"response"`
}

type SchedulerResponse struct {
	Token string `json:"key"`
}

/* Inference Handler: Scheduler Server에 Inference Request 및 Response 전달 */
func (h *Handler) inferHandler(w http.ResponseWriter, r *http.Request) {
	// Extract model information from the URL
	vars := mux.Vars(r)
	model := vars["name"]
	version := vars["version"]

	// Extract request content from the body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("** (ERROR)", err)
		return
	}
	defer r.Body.Close()

	// Triton Inference Request
	response, err := requestScheduler(body, model, version)
	if err != nil {
		log.Println("** (ERROR)", err)
		rend.JSON(w, http.StatusBadRequest, nil)
		return
	}

	// Deliver the response to the client
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

/* Send an inference request to the Scheduler-node and return the request */
func requestScheduler(request []byte, model string, version string) ([]byte, error) {
	log.Println("* (System) Request: ▽▽▽▽▽▽▽▽▽▽")
	log.Println("Model:", model)
	log.Println("Version:", version)
	log.Println(string(request))

	// Search for available ports
	listenPort, err := networkController.GetAvailablePort()
	if err != nil {
		return nil, err
	}

	// Url setting
	urlQuery := "?model=" + model + "&version=" + version + "&address=" + networkController.GetLocalIP() + ":" + strconv.Itoa(listenPort)
	url := "http://" + setting.SchedulerUrl + "/request" + urlQuery

	// Request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(request))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Scheduler Server Response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Start listening when the status code is 200
	if resp.StatusCode == http.StatusOK {
		log.Println("* (System) Received a status code 200 from the Scheduler.")

		// Read body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		// Json decode
		var schedulerResponse SchedulerResponse
		if err := json.Unmarshal(body, &schedulerResponse); err != nil {
			return nil, err
		}
		log.Println("* (System) Get token:", schedulerResponse.Token)

		// Listen and get response
		response, err := listen(listenPort, schedulerResponse.Token)
		if err != nil {
			return nil, err
		}

		return response, nil
	} else {
		return nil, errors.New("The token value is invalid.")
	}
}

/* Listen on the specified port and deliver the response */
func listen(port int, Token string) ([]byte, error) {
	log.Println("* (System) Listen port:", port)

	// TCP listen
	receiverAddr, err := net.ResolveTCPAddr("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	receiverConn, err := net.ListenTCP("tcp", receiverAddr)
	if err != nil {
		return nil, err
	}
	defer receiverConn.Close()

	buffer := make([]byte, 17485760)
	for {
		//TCP accept
		conn, err := receiverConn.AcceptTCP()
		if err != nil {
			return nil, err
		}
		defer conn.Close()

		// Read
		n, err_ := conn.Read(buffer)
		if err_ != nil {
			return nil, err
		}

		// Json decode
		var tritonResponse TritonResponse
		if err := json.Unmarshal(buffer[:n], &tritonResponse); err != nil {
			return nil, err
		}
		log.Println("* (System) Receive token:", tritonResponse.Token)

		// Check the token value
		if tritonResponse.Token == Token {
			log.Println("* (System) Successfully authenticated the token.")
			return []byte(tritonResponse.Response), nil
		} else {
			log.Println("ERROR")

			conn.Close()
		}
	}
}
