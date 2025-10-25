package platform

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, r *http.Request, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func Error(w http.ResponseWriter, r *http.Request, status int, data error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}
