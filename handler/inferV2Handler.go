package handler

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ahr-i/triton-agent/src/httpController"
	"github.com/ahr-i/triton-gateway-client/agentCommunicator"
	"github.com/ahr-i/triton-gateway-client/schedulerCommunicator"
	"github.com/ahr-i/triton-gateway-client/setting"
	"github.com/gorilla/mux"
)

/* Request inference with information on provider, model, version. */
/* If you want to request inference with only the information on model and version, please use the code in 'inferHandler.go'. */
func (h *Handler) inferV2Handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	model := vars["model"]
	version := vars["version"]

	// Extract request content from the body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("** (ERROR)", err)
		return
	}
	defer r.Body.Close()

	// Get GPU-node address
	var address string
	if setting.SchedulerActive {
		address, err = schedulerCommunicator.GetAgentAddressV2(provider, model, version)
		if err != nil {
			log.Println("** (ERROR)", err)
			rend.JSON(w, http.StatusBadRequest, nil)
			return
		}
	} else {
		address = setting.AgentURL
	}

	// Triton inference request
	response, err := agentCommunicator.InferenceV2(address, provider, model, version, body)
	if err != nil {
		log.Println("** (ERROR)", err)
		rend.JSON(w, http.StatusBadRequest, nil)
		return
	}

	// Deliver the response to the client
	httpController.JSON(w, http.StatusOK, response)
}
