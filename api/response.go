package api

import (
	"encoding/json"
	"net/http"
)

func (a *Handler) jsonResponse(w http.ResponseWriter, v interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		a.log.Errorf("Could not respond with JSON: %v", err)
	}
}

func (a *Handler) emptyResponse(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}
