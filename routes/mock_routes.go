package routes

import (
	"api/handlers"
	"database/sql"

	"github.com/gorilla/mux"
)

func MockRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/api/stats/overview", handlers.GenerateOverviewStats()).Methods("GET")
	router.HandleFunc("/api/data/hourly", handlers.GenerateHourlyData()).Methods("GET")
	router.HandleFunc("/api/data/daily", handlers.GenerateDailyData()).Methods("GET")
	router.HandleFunc("/api/data/weekly", handlers.GenerateWeeklyData()).Methods("GET")
	router.HandleFunc("/api/predictions/hourly", handlers.GenerateHourlyPredictions()).Methods("GET")
	router.HandleFunc("/api/predictions/daily", handlers.GenerateDailyPredictions()).Methods("GET")
	router.HandleFunc("/api/predictions/weekly", handlers.GenerateWeeklyPredictions()).Methods("GET")
	router.HandleFunc("/api/anomalies", handlers.GenerateAnomalies()).Methods("GET")

}
