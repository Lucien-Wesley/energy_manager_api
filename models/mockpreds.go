package models

type HourlyPrediction struct {
	Hour       int     `json:"hour"`
	Actual     float64 `json:"actual"`
	Prediction float64 `json:"prediction"`
	Anomaly    bool    `json:"anomaly"`
}

type DailyPrediction struct {
	Date       string  `json:"date"`
	Actual     float64 `json:"actual"`
	Prediction float64 `json:"prediction"`
}

type WeeklyPrediction struct {
	Week       string  `json:"week"`
	Actual     float64 `json:"actual"`
	Prediction float64 `json:"prediction"`
}
