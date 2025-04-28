package models

type OverviewStats struct {
	CurrentUsage   float64 `json:"currentUsage"`
	PreviousUsage  float64 `json:"previousUsage"`
	PercentChange  float64 `json:"percentChange"`
	PredictedUsage float64 `json:"predictedUsage"`
	AnomaliesCount int     `json:"anomaliesCount"`
	Efficiency     int     `json:"efficiency"`
	Unit           string  `json:"unit"`
}
