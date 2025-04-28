package routes

import (
	"api/handlers"
	"database/sql"

	"github.com/gorilla/mux"
)

func HabitatRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/habitats", handlers.CreateHabitat(db)).Methods("POST")
	router.HandleFunc("/habitats/{user_id}", handlers.GetHabitatsByUserID(db)).Methods("GET")
	router.HandleFunc("/habitats/{adresse}", handlers.GetHabitatsByAdress(db)).Methods("GET")
}

func GestionRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/gestionhabitat", handlers.AddGestionHabitat(db)).Methods("POST")
	router.HandleFunc("/gestionhabitat/{user_id}", handlers.GetHabitatGestionByUserID(db)).Methods("GET")
}

func RoomRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/rooms", handlers.AddRoom(db)).Methods("POST")
	router.HandleFunc("/rooms/{habitat_id}", handlers.GetRoomsByHabitatID(db)).Methods("GET")
}
func ApplianceRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/appliances", handlers.AddAppliance(db)).Methods("POST")
	router.HandleFunc("/appliances/{room_id}", handlers.GetAppliancesByRoomID(db)).Methods("GET")
}
