package models

import "time"

type PredictionsTable struct {
	ID      int       `json:"id"`
	User    int       `json:"user"`
	Habitat int       `json:"habitat"`
	Period  string    `json:"period"`
	Date    time.Time `json:"date"`
}

type PredictionsData struct {
	ID           int     `json:"id"`
	Consommation float64 `json:"consommation"`
	Prediction   int     `json:"prediction"`
}

type AppliancesTable struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Addr string `json:"addr"`
	Port int    `json:"port"`
	Room int    `json:"room"`
}

type ExchangeTable struct {
	ID       int             `json:"id"`
	Type     int             `json:"type"`
	Sender   int             `json:"sender"`
	Messages []MessagesTable `json:"messages"`
}

type MessagesTable struct {
	ID       int       `json:"id"`
	Exchange int       `json:"exchange"`
	Message  string    `json:"message"`
	Date     time.Time `json:"date"`
}

type AnomaliesTable struct {
	ID        int       `json:"id"`
	Date      time.Time `json:"date"`
	Type      int       `json:"type"`
	Habitat   int       `json:"habitat"`
	Appliance int       `json:"appliance"`
}
