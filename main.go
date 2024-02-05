package main

import (
	"log"
	"net/http"

	"github.com/ahr-i/triton-gateway-client/handler"
	"github.com/ahr-i/triton-gateway-client/setting"
	"github.com/ahr-i/triton-gateway-client/src/corsController"
	"github.com/urfave/negroni"
)

/* Main */
func main() {
	mux := handler.CreateHandler()
	handler := negroni.Classic()

	defer mux.Close()

	handler.Use(corsController.SetCors("*", "GET, POST, PUT, DELETE", "*", true))
	handler.UseHandler(mux)

	// HTTP Server Start
	log.Println("* (System) HTTP server start.")
	http.ListenAndServe(":"+setting.ServerPort, handler)
}
