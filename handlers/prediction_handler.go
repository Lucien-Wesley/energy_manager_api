package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"api/functions"
	"api/models"
)

// PredictionHandler handles prediction-related requests
func ConsumptionPrediction(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pred models.PredictionsTable
		if err := json.NewDecoder(r.Body).Decode(&pred); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		query := "SELECT COUNT(*) FROM consumption WHERE habitat_id = ?"
		var count int
		if err := db.QueryRow(query, pred.Habitat).Scan(&count); err != nil {
			http.Error(w, "Database query error", http.StatusInternalServerError)
			return
		}
		if count == 0 {
			http.Error(w, "No new consumption data available", http.StatusBadRequest)
			return
		}

		// 1. Récupérer la première date pour ce habitat_id
		var firstDate time.Time
		err := db.QueryRow(`
			SELECT MIN(date) FROM consumption WHERE habitat_id = $1
		`, pred.Habitat).Scan(&firstDate)

		if err != nil {
			http.Error(w, "Error getting first date", http.StatusInternalServerError)
			return
		}

		// 2. Calculer startDate et endDate
		startDate := firstDate.Truncate(time.Hour)
		endDate := time.Now().Truncate(time.Hour).Add(time.Hour) // heure suivante

		// 3. Requête principale
		consumptionQuery := `
			WITH heures AS (
				SELECT generate_series(
					$2::timestamp,
					$3::timestamp,
					'1 hour'::interval
				) AS heure
			)
			SELECT 
				h.heure,
				COALESCE(SUM(c.consommation), 0) AS total_consommation
			FROM 
				heures h
			LEFT JOIN 
				consumption c
				ON c.date > h.heure - interval '1 hour'
				AND c.date <= h.heure
				AND c.habitat_id = $1
			GROUP BY 
				h.heure
			ORDER BY 
				h.heure;
		`

		rows, err := db.Query(consumptionQuery, pred.Habitat, startDate, endDate)
		if err != nil {
			http.Error(w, "Database query error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var consumptionData []models.PredictionsData
		for rows.Next() {
			var data models.PredictionsData
			if err := rows.Scan(&data.Consommation); err != nil {
				http.Error(w, "Error scanning data", http.StatusInternalServerError)
				return
			}
			consumptionData = append(consumptionData, data)
		}


		var trainData []float64
		for _, data := range consumptionData {
			trainData = append(trainData, data.Consommation)
		}

		predictionResults := functions.Predictions(trainData)
		if predictionResults == nil {
			http.Error(w, "Error generating predictions", http.StatusInternalServerError)
			return
		}

		preds := predictionResults[pred.Period]
		if len(preds) == 0 {
			http.Error(w, "No predictions available for selected period", http.StatusInternalServerError)
			return
		}

		// Insérer dans predictions
		insertQuery := "INSERT INTO predictions (user_id, habitat_id, period, date) VALUES (?, ?, ?, ?) RETURNING id"
		var predictionID int
		err = db.QueryRow(insertQuery, pred.User, pred.Habitat, pred.Period, pred.Date).Scan(&predictionID)
		if err != nil {
			http.Error(w, "Error storing prediction metadata", http.StatusInternalServerError)
			return
		}

		// Insérer chaque valeur dans predictions_datas
		for _, value := range preds {
			_, err = db.Exec("INSERT INTO predictions_datas (consommation, prediction_id) VALUES (?, ?)", value, predictionID)
			if err != nil {
				http.Error(w, "Error storing prediction data", http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":    "Prediction stored successfully",
			"prediction": preds,
		})
	}
}

func GetPredictions(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		period := r.URL.Query().Get("period")

		if userID == "" || period == "" {
			http.Error(w, "Missing user_id or period", http.StatusBadRequest)
			return
		}

		// Détermine combien de prédictions il faut récupérer
		var limit int
		switch period {
		case "1d":
			limit = 24 // 24 heures
		case "1w":
			limit = 24 * 7 // 7 jours
		case "1m":
			limit = 24 * 30 // 30 jours
		default:
			http.Error(w, "Invalid period. Use '1d', '1w' or '1m'", http.StatusBadRequest)
			return
		}

		query := `
			SELECT pd.consommation 
			FROM predictions p
			JOIN predictions_datas pd ON pd.prediction_id = p.id
			WHERE p.user_id = ? AND p.period = ?
			ORDER BY p.date DESC
			LIMIT ?;
		`

		rows, err := db.Query(query, userID, period, limit)
		if err != nil {
			http.Error(w, "Database query error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var predictions []float64
		for rows.Next() {
			var value float64
			if err := rows.Scan(&value); err != nil {
				http.Error(w, "Error reading prediction data", http.StatusInternalServerError)
				return
			}
			predictions = append(predictions, value)
		}

		if len(predictions) == 0 {
			http.Error(w, "No predictions found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"predictions": predictions,
		})
	}
}
