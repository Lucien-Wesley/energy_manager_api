package models

type HourlyData struct {
	Hour  int     `json:"hour"`
	Value float64 `json:"value"`
}

type DailyData struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

type WeeklyData struct {
	Week  string  `json:"week"`
	Value float64 `json:"value"`
}
