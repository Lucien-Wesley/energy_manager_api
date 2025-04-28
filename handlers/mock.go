package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"api/models"
)

func randFloat(min, max float64) float64 {
	return math.Round((min+rand.Float64()*(max-min))*100) / 100
}

// Overview Stats
func GenerateOverviewStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rand.Seed(time.Now().UnixNano())
		stats := models.OverviewStats{
			CurrentUsage:   math.Round(randFloat(10, 15)*100) / 100,
			PreviousUsage:  math.Round(randFloat(12, 17)*100) / 100,
			PercentChange:  math.Round(randFloat(-10, 10)*100) / 100,
			PredictedUsage: math.Round(randFloat(9, 14)*100) / 100,
			AnomaliesCount: rand.Intn(5),
			Efficiency:     rand.Intn(41) + 60,
			Unit:           "kWh",
		}
		json.NewEncoder(w).Encode(stats)
	}
}

// Hourly Data
func GenerateHourlyData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := make([]models.HourlyData, 24)
		for i := 0; i < 24; i++ {
			data[i] = models.HourlyData{
				Hour:  i,
				Value: randFloat(0.5, 2.5),
			}
		}
		json.NewEncoder(w).Encode(data)
	}
}
func getPastDates(days int) []string {
	dates := make([]string, days)
	now := time.Now()

	for i := 0; i < days; i++ {
		date := now.AddDate(0, 0, -(days - 1 - i))
		dates[i] = date.Format(time.RFC3339) // équivalent de toISOString()
	}

	return dates
}

// generateDailyData génère des données de consommation journalière
func GenerateDailyData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dates := getPastDates(7)
		data := make([]models.DailyData, len(dates))
		for i, date := range dates {
			val := math.Round((rand.Float64()*10+15)*100) / 100
			data[i] = models.DailyData{
				Date:  date,
				Value: val,
			}
		}
		json.NewEncoder(w).Encode(data)
	}
}

// generateWeeklyData génère des données de consommation hebdomadaire
func GenerateWeeklyData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := make([]models.WeeklyData, 4)
		for i := 0; i < 4; i++ {
			val := math.Round((rand.Float64()*50+75)*100) / 100
			data[i] = models.WeeklyData{
				Week:  fmt.Sprintf("Week %d", i+1),
				Value: val,
			}
		}
		json.NewEncoder(w).Encode(data)
	}
}

// Hourly Predictions
func GenerateHourlyPredictions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := make([]models.HourlyPrediction, 24)
		for i := 0; i < 24; i++ {
			actual := randFloat(0.5, 2.5)
			variation := randFloat(-0.15, 0.15)
			predicted := actual * (1 - variation)
			data[i] = models.HourlyPrediction{
				Hour:       i,
				Actual:     math.Round(actual*100) / 100,
				Prediction: math.Round(predicted*100) / 100,
				Anomaly:    math.Abs(variation) > 0.12,
			}
		}
		json.NewEncoder(w).Encode(data)
	}
}

// Daily Predictions (7 days)
func GenerateDailyPredictions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := []models.DailyPrediction{}
		for i := 0; i < 7; i++ {
			date := time.Now().AddDate(0, 0, -6+i).Format(time.RFC3339)
			actual := randFloat(15, 25)
			pred := actual * (0.9 + rand.Float64()*0.2)
			data = append(data, models.DailyPrediction{
				Date:       date,
				Actual:     math.Round(actual*100) / 100,
				Prediction: math.Round(pred*100) / 100,
			})
		}
		json.NewEncoder(w).Encode(data)
	}
}

// Weekly Predictions (4 weeks)
func GenerateWeeklyPredictions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := []models.WeeklyPrediction{}
		for i := 0; i < 4; i++ {
			actual := randFloat(75, 125)
			pred := actual * (0.9 + rand.Float64()*0.2)
			data = append(data, models.WeeklyPrediction{
				Week:       "Week " + strconv.Itoa(i+1),
				Actual:     actual,
				Prediction: pred,
			})
		}
		json.NewEncoder(w).Encode(data)
	}
}

// Anomalies
func GenerateAnomalies() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rooms := []string{"Living Room", "Kitchen", "Bedroom", "Bathroom", "Home Office"}
		devices := []string{"TV", "Refrigerator", "Air Conditioner", "Lights", "Computer", "Microwave", "Dishwasher", "Water Heater", "Washing Machine"}
		severities := []string{"low", "medium", "high"}

		anomalies := []models.Anomaly{}
		for i := 0; i < 6; i++ {
			expected := randFloat(0.5, 2.5)
			dev := randFloat(10, 50)
			actual := expected * (1 + dev/100)
			severity := severities[0]
			if dev > 35 {
				severity = "high"
			} else if dev > 20 {
				severity = "medium"
			}
			anomalies = append(anomalies, models.Anomaly{
				ID:            "anomaly-" + strconv.Itoa(i),
				Timestamp:     time.Now().Add(-time.Duration(rand.Intn(72)) * time.Hour).Format(time.RFC3339),
				Room:          rooms[rand.Intn(len(rooms))],
				Device:        devices[rand.Intn(len(devices))],
				ExpectedValue: expected,
				ActualValue:   actual,
				Deviation:     dev,
				Severity:      severity,
			})
		}
		json.NewEncoder(w).Encode(anomalies)
	}
}
