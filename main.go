package main

import (
	"log"
	"net/http"

	"api/config"
	"api/middleware"
	"api/routes"

	"github.com/gorilla/mux"
)

func main() {
	db := config.ConnectDB()
	defer db.Close()

	router := mux.NewRouter()
	router.Use(middleware.JSONMiddleware)

	routes.UserRoutes(router, db)
	routes.HabitatRoutes(router, db)
	routes.GestionRoutes(router, db)
	routes.RoomRoutes(router, db)
	routes.ApplianceRoutes(router, db)
	routes.ConsumptionRoutes(router, db)
	routes.ExchangeRoutes(router, db)
	routes.PredictionRoutes(router, db)
	routes.TestESP32Routes(router)
	routes.MockRoutes(router, db)

	log.Println("Serveur en écoute sur le port 5400...")
	log.Fatal(http.ListenAndServe(":5400", router))
}

//F^655377299994aj
