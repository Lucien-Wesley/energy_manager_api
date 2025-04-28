package routes

import (
	"api/handlers"
	"database/sql"

	"github.com/gorilla/mux"
)

func PredictionRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/predict", handlers.ConsumptionPrediction(db)).Methods("POST")
	router.HandleFunc("/predictions", handlers.GetPredictions(db)).Methods("GET")
}
