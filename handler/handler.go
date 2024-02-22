package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var rend *render.Render = render.New()

type Handler struct {
	http.Handler
}

func CreateHandler() *Handler {
	mux := mux.NewRouter()
	handler := &Handler{
		Handler: mux,
	}

	mux.HandleFunc("/", handler.homeHandler).Methods("GET")                                                                               // HTML, CSS, JS 요청
	mux.HandleFunc("/ping", handler.pingHandler).Methods("GET")                                                                           // Ping Check
	mux.HandleFunc("/model/{model:[a-z-_]+}/{version:[0-9]+}/infer", handler.inferHandler).Methods("POST")                                // Scheduler Server Inference 요청
	mux.HandleFunc("/provider/{provider:[a-z-_]+}/model/{model:[a-z-_]+}/{version:[0-9]+}/infer", handler.inferV2Handler).Methods("POST") // Request inference with information on provider, model, version.

	return handler
}
