package main

import (
	"net/http"

	"github.com/ahr-i/triton-gateway-client/handler"
	"github.com/ahr-i/triton-gateway-client/models"
	"github.com/ahr-i/triton-gateway-client/setting"
	"github.com/ahr-i/triton-gateway-client/src/corsController"
	"github.com/urfave/negroni"
)

/* Server Setting */
func Init() {
	models.Init(setting.ModelPath) // Model List Init
}

/* Main */
func main() {
	Init()

	mux := handler.CreateHandler()
	handler := negroni.Classic()

	defer mux.Close()

	handler.Use(corsController.SetCors("*", "GET, POST, PUT, DELETE", "*", true))
	handler.UseHandler(mux)

	// HTTP Server Start
	http.ListenAndServe(":"+setting.ServerPort, handler)
}
