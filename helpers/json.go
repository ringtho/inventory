package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSON is a helper function to write JSON responses
func JSON(w http.ResponseWriter, status int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

// respondWithError is a helper function to write error responses
func RespondWithError(w http.ResponseWriter, status int, message string){

	type errorResponse struct {
		Error string `json:"error"`
	}

	JSON(w, status, errorResponse{Error: message})
}