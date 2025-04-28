package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"api/models"

	"github.com/gorilla/mux"
)

func AddApplianceConsumption(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ac models.ConsommationTable
		if err := json.NewDecoder(r.Body).Decode(&ac); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Insérer la consommation de l'appareil
		err := db.QueryRow("INSERT INTO appliance_consumption (appliance_id, date, consumption) VALUES ($1, $2, $3) RETURNING id",
			ac.Appliance, time.Now(), ac.Consumption).Scan(&ac.ID)
		if err != nil {
			log.Println("Erreur lors de la création de la consommation de l'appareil :", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(ac)
	}
}

func GetApplianceConsumption(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		applianceID := vars["appl_id"]
		period := r.URL.Query().Get("period") // "daily", "weekly", or "monthly"

		var query string
		switch period {
		case "daily":
			query = `
				SELECT DATE(date) AS period, SUM(consumption) AS total_consumption
				FROM appliance_consumption
				WHERE appliance_id = $1
				GROUP BY DATE(date)
				ORDER BY DATE(date) DESC`
		case "weekly":
			query = `
				SELECT DATE_TRUNC('week', date) AS period, SUM(consumption) AS total_consumption
				FROM appliance_consumption
				WHERE appliance_id = $1
				GROUP BY DATE_TRUNC('week', date)
				ORDER BY DATE_TRUNC('week', date) DESC`
		case "monthly":
			query = `
				SELECT DATE_TRUNC('month', date) AS period, SUM(consumption) AS total_consumption
				FROM appliance_consumption
				WHERE appliance_id = $1
				GROUP BY DATE_TRUNC('month', date)
				ORDER BY DATE_TRUNC('month', date) DESC`
		default:
			http.Error(w, "Invalid or missing period parameter. Use 'daily', 'weekly', or 'monthly'.", http.StatusBadRequest)
			return
		}

		rows, err := db.Query(query, applianceID)
		if err != nil {
			log.Println("Error querying appliance consumption:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var results []map[string]interface{}
		for rows.Next() {
			var period string
			var totalConsumption float64
			if err := rows.Scan(&period, &totalConsumption); err != nil {
				log.Println("Error scanning row:", err)
				http.Error(w, "Server error", http.StatusInternalServerError)
				return
			}
			results = append(results, map[string]interface{}{
				"period":            period,
				"total_consumption": totalConsumption,
			})
		}

		if err := rows.Err(); err != nil {
			log.Println("Error iterating rows:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}

func AddHabitatConsumption(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var hc models.ConsommationsTable
		if err := json.NewDecoder(r.Body).Decode(&hc); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Insérer la consommation de l'habitat
		err := db.QueryRow("INSERT INTO habitat_consumption (habitat_id, date, consumption) VALUES ($1, $2, $3) RETURNING id",
			hc.Habitat, time.Now(), hc.Consumption).Scan(&hc.ID)
		if err != nil {
			log.Println("Erreur lors de la création de la consommation de l'habitat :", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(hc)
	}
}

func GetHabitatConsumption(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		habitatID := vars["habitat_id"]
		period := r.URL.Query().Get("period") // "daily", "weekly", or "monthly"

		var query string
		switch period {
		case "daily":
			query = `
				SELECT DATE(date) AS period, SUM(consumption) AS total_consumption
				FROM habitat_consumption
				WHERE habitat_id = $1
				GROUP BY DATE(date)
				ORDER BY DATE(date) DESC`
		case "weekly":
			query = `
				SELECT DATE_TRUNC('week', date) AS period, SUM(consumption) AS total_consumption
				FROM habitat_consumption
				WHERE habitat_id = $1
				GROUP BY DATE_TRUNC('week', date)
				ORDER BY DATE_TRUNC('week', date) DESC`
		case "monthly":
			query = `
				SELECT DATE_TRUNC('month', date) AS period, SUM(consumption) AS total_consumption
				FROM habitat_consumption
				WHERE habitat_id = $1
				GROUP BY DATE_TRUNC('month', date)
				ORDER BY DATE_TRUNC('month', date) DESC`
		default:
			http.Error(w, "Invalid or missing period parameter. Use 'daily', 'weekly', or 'monthly'.", http.StatusBadRequest)
			return
		}

		rows, err := db.Query(query, habitatID)
		if err != nil {
			log.Println("Error querying habitat consumption:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var results []map[string]interface{}
		for rows.Next() {
			var period string
			var totalConsumption float64
			if err := rows.Scan(&period, &totalConsumption); err != nil {
				log.Println("Error scanning row:", err)
				http.Error(w, "Server error", http.StatusInternalServerError)
				return
			}
			results = append(results, map[string]interface{}{
				"period":            period,
				"total_consumption": totalConsumption,
			})
		}

		if err := rows.Err(); err != nil {
			log.Println("Error iterating rows:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}
