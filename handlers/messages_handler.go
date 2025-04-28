package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"api/models"

	"github.com/gorilla/mux"
)

func CreateMessage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var exchange models.ExchangeTable
		if err := json.NewDecoder(r.Body).Decode(&exchange); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Insert the exchange into the ExchangeTable
		exchangeQuery := "INSERT INTO exchanges (type, sender) VALUES ($1, $2) RETURNING ID"
		err := db.QueryRow(exchangeQuery, exchange.Type, exchange.Sender).Scan(&exchange.ID)
		if err != nil {
			http.Error(w, "Failed to create exchange", http.StatusInternalServerError)
			return
		}

		var message models.MessagesTable
		if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
			http.Error(w, "Invalid request payload for message", http.StatusBadRequest)
			return
		}

		// Assign the generated exchange ID to the message
		message.Exchange = exchange.ID

		// Insert the message into the MessagesTable
		messageQuery := "INSERT INTO messages (exchange, message, date) VALUES ($1, $2, $3) RETURNING ID"
		err = db.QueryRow(messageQuery, message.Exchange, message.Message, message.Date).Scan(&message.ID)
		if err != nil {
			http.Error(w, "Failed to create message", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		exchange.Messages = append(exchange.Messages, message)
		json.NewEncoder(w).Encode(exchange)
	}
}	

func GetExchangesMessagesBySenderID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		senderID := vars["sender_id"]

		rows, err := db.Query("SELECT e.id, e.type, e.sender, m.id, m.message, m.date FROM exchanges e JOIN messages m ON e.id = m.exchange WHERE e.sender = $1", senderID)
		if err != nil {
			http.Error(w, "Error fetching exchanges and messages", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var exchanges []models.ExchangeTable
		for rows.Next() {
			var exchange models.ExchangeTable
			var message models.MessagesTable
			if err := rows.Scan(&exchange.ID, &exchange.Type, &exchange.Sender, &message.ID, &message.Message, &message.Date); err != nil {
				http.Error(w, "Error scanning exchange and message", http.StatusInternalServerError)
				return
			}
			exchange.Messages = append(exchange.Messages, message)
			exchanges = append(exchanges, exchange)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Error iterating over exchanges and messages", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(exchanges)
	}
}