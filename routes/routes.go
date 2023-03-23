package routes

import (
	"net/http"

	"github.com/Rishabh/golang-ws/controllers"
	"github.com/bmizerany/pat"
)

func Routes() http.Handler {
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(controllers.Home))
	mux.Get("/ws", http.HandlerFunc(controllers.WsEndpoint))
	
	return mux
}
