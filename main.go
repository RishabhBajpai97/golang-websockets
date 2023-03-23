package main

import (
	"net/http"

	"github.com/Rishabh/golang-ws/controllers"
	"github.com/Rishabh/golang-ws/routes"
)

func main() {
	mux := routes.Routes()
		go controllers.ListentoWs()
	_ = http.ListenAndServe(":8080", mux)
}
