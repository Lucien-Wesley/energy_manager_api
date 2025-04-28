package models

import "time"

type ConsommationTable struct {
	ID           int       `json:"id"`
	Date         time.Time `json:"date"`
	Consumption float64   `json:"consommation"`
	Appliance    int       `json:"appliance"`
}

type ConsommationsTable struct {
	ID           int       `json:"id"`
	Date         time.Time `json:"date"`
	Consumption float64   `json:"consommation"`
	Habitat         int       `json:"home"`
}
