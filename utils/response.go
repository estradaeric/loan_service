package utils

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// WriteSuccess sends a standardized success response
func WriteSuccess(w http.ResponseWriter, status int, message string, data interface{}) {
	resp := APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
	writeJSON(w, status, resp)
}

// WriteError sends a standardized error response
func WriteError(w http.ResponseWriter, status int, message string) {
	resp := APIResponse{
		Status:  "error",
		Message: message,
	}
	writeJSON(w, status, resp)
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}