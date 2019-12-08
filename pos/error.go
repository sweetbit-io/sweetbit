package pos

import (
	"encoding/json"
	"net/http"
)

type errorMessage struct {
	Error string `json:"error"`
}

func (p *Handler) jsonError(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(&errorMessage{
		Error: error,
	})
	if err != nil {
		p.log.Errorf("Could not respond with error: %v", err)
	}
}
