package agentCommunicator

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

/* Send an inference request to the agent. */
/* Information on provider, model, and version is required. */
func InferenceV2(address string, provider string, model string, version string, request []byte) ([]byte, error) {
	log.Println("* (System) Send a request to the Triton agent.")
	url := fmt.Sprintf("http://%s/provider/%s/model/%s/%s/infer", address, provider, model, version)

	// Send an inference request to the agent.
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(request))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Println("* (System) Successfully received a response from the Triton agent.")

	return body, nil
}

/* Send an inference request to the agent. */
func Inference(address string, model string, version string, request []byte) ([]byte, error) {
	log.Println("* (System) Send a request to the Triton agent.")
	url := fmt.Sprintf("http://%s/model/%s/%s/infer", address, model, version)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(request))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Println("* (System) Successfully received a response from the Triton agent.")

	return body, nil
}
