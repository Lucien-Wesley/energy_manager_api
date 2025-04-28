package routes

import (
	"api/handlers"
	"database/sql"

	"github.com/gorilla/mux"
)

func UserRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/users", handlers.GetUsers(db)).Methods("GET")
	router.HandleFunc("/users/{id}", handlers.GetUser(db)).Methods("GET")
	router.HandleFunc("/users", handlers.CreateUser(db)).Methods("POST")
	router.HandleFunc("/users/login", handlers.LoginUser(db)).Methods("POST")
}
