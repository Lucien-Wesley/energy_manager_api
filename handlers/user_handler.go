package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"api/models" // Remplace par le bon chemin de ton package models
	"api/utils"

	"github.com/gorilla/mux"
)

func CreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u models.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Hashage du mot de passe
		hashedPassword, err := utils.HashPassword(u.Password)
		//fmt.Println(hashedPassword)
		if err != nil {
			http.Error(w, "Erreur lors du hashage du mot de passe", http.StatusInternalServerError)
			return
		}

		// Insérer l'utilisateur avec le mot de passe haché
		err = db.QueryRow("INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id",
			u.Name, u.Email, hashedPassword).Scan(&u.ID)
		if err != nil {
			log.Println("Erreur lors de la création de l'utilisateur :", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(u)
	}
}

func GetUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var u models.User
		err := db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(u)
	}
}

func GetUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, email, password FROM users")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		users := []models.User{}
		for rows.Next() {
			var u models.User
			if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password); err != nil {
				log.Fatal(err)
			}
			users = append(users, u)
		}

		json.NewEncoder(w).Encode(users)
	}
}

func LoginUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		json.NewDecoder(r.Body).Decode(&u)

		var storedPassword string
		var userID int

		// Vérifier si l'utilisateur existe
		err := db.QueryRow("SELECT id, password FROM users WHERE email = $1", u.Email).Scan(&userID, &storedPassword)
		if err != nil {
			http.Error(w, "Utilisateur introuvable", http.StatusUnauthorized)
			return
		}

		// Vérification du mot de passe
		if !utils.CheckPassword(storedPassword, u.Password) {
			http.Error(w, "Mot de passe incorrect", http.StatusUnauthorized)
			return
		}

		// Générer un token
		token, err := utils.GenerateToken(userID)
		if err != nil {
			http.Error(w, "Erreur interne", http.StatusInternalServerError)
			return
		}

		// Répondre avec le token
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token, "user_id": strconv.Itoa(userID)})
	}
}
