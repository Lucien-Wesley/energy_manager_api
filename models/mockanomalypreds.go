package models

type Anomaly struct {
	ID            string  `json:"id"`
	Timestamp     string  `json:"timestamp"`
	Room          string  `json:"room"`
	Device        string  `json:"device"`
	ExpectedValue float64 `json:"expectedValue"`
	ActualValue   float64 `json:"actualValue"`
	Deviation     float64 `json:"deviation"`
	Severity      string  `json:"severity"`
}
