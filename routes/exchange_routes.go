package routes

import (
	"api/handlers"
	"database/sql"

	"github.com/gorilla/mux"
)

func ExchangeRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/exchanges", handlers.CreateMessage(db)).Methods("POST")
	router.HandleFunc("/exchanges/{sender_id}", handlers.GetExchangesMessagesBySenderID(db)).Methods("GET")
}
