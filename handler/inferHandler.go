package handler

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ahr-i/triton-agent/src/httpController"
	"github.com/ahr-i/triton-gateway-client/agentCommunicator"
	"github.com/ahr-i/triton-gateway-client/schedulerCommunicator"
	"github.com/gorilla/mux"
)

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

	// Get GPU-node address
	address, err := schedulerCommunicator.GetAgentAddress(model, version)
	if err != nil {
		log.Println("** (ERROR)", err)
		rend.JSON(w, http.StatusBadRequest, nil)
		return
	}

	// Triton inference request
	response, err := agentCommunicator.Inference(address, model, version, body)
	if err != nil {
		log.Println("** (ERROR)", err)
		rend.JSON(w, http.StatusBadRequest, nil)
		return
	}

	// Deliver the response to the client
	httpController.JSON(w, http.StatusOK, response)
}
