package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"api/models"

	"github.com/gorilla/mux"
)

func CreateHabitat(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var h models.Habitat
		if err := json.NewDecoder(r.Body).Decode(&h); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Insérer l'habitat
		err := db.QueryRow("INSERT INTO habitats (adresse, user) VALUES ($1, $2) RETURNING id",
			h.Adresse, h.User).Scan(&h.ID)
		if err != nil {
			log.Println("Erreur lors de la création de l'habitat :", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(h)
	}
}

func GetHabitatsByUserID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["user_id"]

		rows, err := db.Query("SELECT adresse FROM habitats WHERE user = $1", userID)
		if err != nil {
			http.Error(w, "Error fetching habitats", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var addresses []string
		for rows.Next() {
			var adresse string
			if err := rows.Scan(&adresse); err != nil {
				http.Error(w, "Error scanning habitat", http.StatusInternalServerError)
				return
			}
			addresses = append(addresses, adresse)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Error iterating over habitats", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(addresses)
	}
}

func GetHabitatsByAdress(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		adresse := vars["adresse"]

		rows, err := db.Query("SELECT id FROM habitats WHERE user = $1", adresse)
		if err != nil {
			http.Error(w, "Error fetching habitats", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var id int
		if rows.Next() {
			if err := rows.Scan(&id); err != nil {
				http.Error(w, "Error scanning habitat", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "No habitat found for the given address", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(id)
	}
}

func AddGestionHabitat(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var h models.HabitatGestion
		if err := json.NewDecoder(r.Body).Decode(&h); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Insérer l'habitat
		err := db.QueryRow("INSERT INTO habitat_gestion (user_id, habitat_id, validation) VALUES ($1, $2, $3) RETURNING id",
			h.User, h.Habitat, h.Validation).Scan(&h.ID)
		if err != nil {
			log.Println("Erreur lors de la création de l'habitat :", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(h)
	}
}

func GetHabitatGestionByUserID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["user_id"]

		rows, err := db.Query("SELECT habitat FROM habitat_gestion WHERE user_id = $1", userID)
		if err != nil {
			http.Error(w, "Error fetching habitats", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var addresses []string
		for rows.Next() {
			var adresse string
			if err := rows.Scan(&adresse); err != nil {
				http.Error(w, "Error scanning habitat", http.StatusInternalServerError)
				return
			}
			addresses = append(addresses, adresse)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Error iterating over habitats", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(addresses)
	}
}

func AddRoom(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var room models.RoomsTable
		if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Insérer la salle
		err := db.QueryRow("INSERT INTO rooms (name, habitat_id) VALUES ($1, $2) RETURNING id",
			room.Name, room.Habitat).Scan(&room.ID)
		if err != nil {
			log.Println("Erreur lors de la création de la salle :", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(room)
	}
}

func GetRoomsByHabitatID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		habitatID := vars["habitat_id"]

		rows, err := db.Query("SELECT id, name FROM rooms WHERE habitat_id = $1", habitatID)
		if err != nil {
			http.Error(w, "Error fetching rooms", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var rooms []map[string]interface{}
		for rows.Next() {
			var id int
			var name string
			if err := rows.Scan(&id, &name); err != nil {
				http.Error(w, "Error scanning room", http.StatusInternalServerError)
				return
			}
			rooms = append(rooms, map[string]interface{}{
				"id":   id,
				"name": name,
			})
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Error iterating over rooms", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(rooms)
	}
}

func AddAppliance(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var appliance models.AppliancesTable
		if err := json.NewDecoder(r.Body).Decode(&appliance); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Insérer l'appareil
		err := db.QueryRow("INSERT INTO appliances (name, addr, port, room_id) VALUES ($1, $2, $3, $4) RETURNING id",
			appliance.Name, appliance.Addr, appliance.Port, appliance.Room).Scan(&appliance.ID)
		if err != nil {
			log.Println("Erreur lors de la création de l'appareil :", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(appliance)
	}
}

func GetAppliancesByRoomID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		roomID := vars["room_id"]

		rows, err := db.Query("SELECT id, name, addr, port, room FROM appliances WHERE room = $1", roomID)
		if err != nil {
			http.Error(w, "Error fetching appliances", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var appliances []map[string]interface{}
		for rows.Next() {
			var id int
			var name, addr string
			var port int
			var room int
			if err := rows.Scan(&id, &name, &addr, &port, &room); err != nil {
				http.Error(w, "Error scanning appliance", http.StatusInternalServerError)
				return
			}
			appliances = append(appliances, map[string]interface{}{
				"id":   id,
				"name": name,
				"addr": addr,
				"port": port,
				"room": room,
			})
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Error iterating over appliances", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(appliances)
	}
}
