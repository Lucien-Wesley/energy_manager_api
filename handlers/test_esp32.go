package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type TestRequest struct {
	Message string `json:"message"`
}

func TestESP32Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req TestRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	response := map[string]string{
		"received": req.Message,
	}

	fmt.Println("Response:", response)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
