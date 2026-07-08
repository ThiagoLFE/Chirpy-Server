package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type ValidateRequest struct {
	Body string `json:"body"`
}

type ValidateResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

func respondWithError(w http.ResponseWriter, statusCode int, msg string) {

	respondWithJSON(w, statusCode, ErrorResponse{Error: msg})
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	data, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}
