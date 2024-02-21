package schedulerCommunicator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ahr-i/triton-gateway-client/setting"
)

/* Scheduler Struct */
type schedulerResponse struct {
	Address string `json:"address"`
}

/* Obtain the address of the GPU node from the scheduler. */
/* Information on provider, model, and version is required. */
func GetAgentAddressV2(provider string, model string, version string) (string, error) {
	log.Println("* (System) Model information: ▽▽▽▽▽▽▽▽▽▽")
	log.Println("provider:", provider)
	log.Println("Model:", model)
	log.Println("Version:", version)

	// Url setting
	urlQuery := fmt.Sprintf("?id=%s&model=%s&version=%s", provider, model, version)
	url := fmt.Sprintf("http://%s/request%s", setting.SchedulerUrl, urlQuery)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return "", err
	}

	// Start listening when the status code is 200
	if resp.StatusCode == http.StatusOK {
		log.Println("* (System) Received a status code 200 from the Scheduler.")

		// Read body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		// Json decode
		var schedulerResponse schedulerResponse
		if err := json.Unmarshal(body, &schedulerResponse); err != nil {
			return "", err
		}
		log.Println("* (System) Get GPU-node address:", schedulerResponse.Address)

		return schedulerResponse.Address, nil
	} else {
		return "", errors.New("scheduler error")
	}
}

/* Obtain the address of the GPU node from the scheduler. */
func GetAgentAddress(model string, version string) (string, error) {
	log.Println("* (System) Model information: ▽▽▽▽▽▽▽▽▽▽")
	log.Println("Model:", model)
	log.Println("Version:", version)

	// Url setting
	urlQuery := "?model=" + model + "&version=" + version
	url := "http://" + setting.SchedulerUrl + "/request" + urlQuery

	// Request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// Scheduler Server Response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Start listening when the status code is 200
	if resp.StatusCode == http.StatusOK {
		log.Println("* (System) Received a status code 200 from the Scheduler.")

		// Read body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		// Json decode
		var schedulerResponse schedulerResponse
		if err := json.Unmarshal(body, &schedulerResponse); err != nil {
			return "", err
		}
		log.Println("* (System) Get GPU-node address:", schedulerResponse.Address)

		return schedulerResponse.Address, nil
	} else {
		return "", errors.New("scheduler error")
	}
}
