package routes

import (
	"api/handlers"

	"github.com/gorilla/mux"
)

func TestESP32Routes(router *mux.Router) {
	router.HandleFunc("/testesp32", handlers.TestESP32Handler).Methods("POST")
}
