package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// ConnectDB initialise la connexion à la base de données
func ConnectDB() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Erreur de connexion à la base de données:", err)
	}

	// Vérifier la connexion
	if err := db.Ping(); err != nil {
		log.Fatal("Impossible de ping la base de données:", err)
	}

	fmt.Println("✅ Connexion réussie à la base de données")
	InitializeTables(db)

	return db
}

// InitializeTables crée les tables si elles n'existent pas encore
func InitializeTables(db *sql.DB) {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(100) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS habitats (
			id SERIAL PRIMARY KEY,
			adresse VARCHAR(255) NOT NULL,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS habitat_gestion (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			habitat_id INTEGER REFERENCES habitats(id),
			validation BOOLEAN DEFAULT false
		)`,
		`CREATE TABLE IF NOT EXISTS predictions (
			id SERIAL PRIMARY KEY,
			period VARCHAR NOT NULL,
			habitat_id INTEGER REFERENCES habitats(id) ON DELETE CASCADE,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			date TIMESTAMP NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS rooms (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			habitat_id INTEGER REFERENCES habitats(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS appliances (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			adresse VARCHAR(100) NOT NULL,
			port INTEGER NOT NULL,
			state BOOLEAN DEFAULT false,
			room_id INTEGER REFERENCES rooms(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS appliance_consumption (
			id SERIAL PRIMARY KEY,
			date TIMESTAMP NOT NULL,
			consumption FLOAT NOT NULL,
			appliance_id INTEGER REFERENCES appliances(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS habitat_consumption (
			id SERIAL PRIMARY KEY,
			date TIMESTAMP NOT NULL,
			consumption FLOAT NOT NULL,
			habitat_id INTEGER REFERENCES habitats(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS exchanges (
			id SERIAL PRIMARY KEY,
			type INTEGER NOT NULL,
			sender_id INTEGER REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS messages (
			id SERIAL PRIMARY KEY,
			exchange_id INTEGER REFERENCES exchanges(id) ON DELETE CASCADE,
			message TEXT NOT NULL,
			date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS anomalies (
			id SERIAL PRIMARY KEY,
			date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			type INTEGER NOT NULL,
			habitat_id INTEGER REFERENCES habitats(id) ON DELETE CASCADE,
			appliance_id INTEGER REFERENCES appliances(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS predictions_datas (
			id SERIAL PRIMARY KEY,
			consommation FLOAT NOT NULL,
			prediction_id INTEGER REFERENCES predictions(id) ON DELETE CASCADE
		)`,
	}

	for _, query := range tables {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatal("❌ Erreur lors de la création des tables:", err)
		}
	}

	fmt.Println("✅ Toutes les tables sont créées ou existent déjà")
}
