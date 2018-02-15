package client

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Msg string `json:"message"`
}

func WriteErrorResponse(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	er := ErrorResponse{msg}
	b, err := json.Marshal(er)
	if err != nil {
		log.Printf("Error Marshaling JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(b)
}
