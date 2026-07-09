package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github/ThiagoLFE/Chirpy-Server/internal/database"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type ChirpCmd struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {

	var chirpCmd ChirpCmd
	if err := json.NewDecoder(r.Body).Decode(&chirpCmd); err != nil {
		respondWithError(w, http.StatusInternalServerError, "fail to decode payload: "+err.Error())
		return
	}

	chirpBody := strings.TrimSpace(chirpCmd.Body)
	if len(chirpBody) == 0 {
		respondWithError(w, http.StatusBadRequest, "body is required")
		return
	}

	if len(chirpCmd.UserID) == 0 {
		respondWithError(w, http.StatusBadRequest, "user is required")
		return
	}

	u, err := cfg.db.GetUserByID(r.Context(), chirpCmd.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusBadGateway, "user not exist")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "fail to get user: "+err.Error())
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{Body: chirpBody, UserID: u.ID})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "fail to create chirp: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}
