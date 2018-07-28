package client

import (
	"encoding/json"
	"net/http"
)

// Error response is a JSON wrapper around an error message
type ErrorResponse struct {
	Msg string `json:"message"`
}

// WriteErrorResponse writes a generic HTTP JSON error
func WriteErrorResponse(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	er := ErrorResponse{msg}
	b, err := json.Marshal(er)
	if err != nil {
		clientLog.Errorf("Error Marshaling JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(b)
}
