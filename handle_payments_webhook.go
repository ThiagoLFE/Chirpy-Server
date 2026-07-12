package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type PaymentWebHookRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlePaymentsWebHook(w http.ResponseWriter, r *http.Request) {
	var paymentWebHookRequest PaymentWebHookRequest

	if err := json.NewDecoder(r.Body).Decode(&paymentWebHookRequest); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid body")
		return
	}

	if paymentWebHookRequest.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if _, err := cfg.db.UpdateUserRedChirpMember(r.Context(), paymentWebHookRequest.Data.UserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Printf("fail to apply red chirp member to user %s: %v", paymentWebHookRequest.Data.UserID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
