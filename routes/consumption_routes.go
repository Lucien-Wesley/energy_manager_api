package routes

import (
	"api/handlers"
	"database/sql"

	"github.com/gorilla/mux"
)

func ConsumptionRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/consumption", handlers.AddApplianceConsumption(db)).Methods("POST")
	router.HandleFunc("/consumption/{user_id}", handlers.GetApplianceConsumption(db)).Methods("GET")
	router.HandleFunc("/consumptions", handlers.AddHabitatConsumption(db)).Methods("POST")
	router.HandleFunc("/consumption/{user_id}", handlers.GetHabitatConsumption(db)).Methods("GET")
}