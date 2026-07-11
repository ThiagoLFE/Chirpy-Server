package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github/ThiagoLFE/Chirpy-Server/internal/database"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ChirpCmd struct {
	Body string `json:"body"`
}

type ChirpResponse struct {
	ID        uuid.UUID `json:"id"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const MaxChirpBodyLength = 140

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {

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
	if len(chirpBody) > MaxChirpBodyLength {
		respondWithError(w, http.StatusBadRequest, "body is too large, max of 140 characters")
		return
	}

	userID, ok := r.Context().Value(userIDContextKey).(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	dbChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{Body: chirpBody, UserID: userID})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "fail to create chirp: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, ChirpResponse{
		ID:        dbChirp.ID,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
	})
}

func (cfg *apiConfig) handleListChirps(w http.ResponseWriter, r *http.Request) {
	dbChips, err := cfg.db.ListChirps(r.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithJSON(w, http.StatusOK, []map[string]string{})
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	formattedList := make([]ChirpResponse, 0)
	for _, chirp := range dbChips {
		formattedList = append(formattedList, ChirpResponse{
			ID:        chirp.ID,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, formattedList)
}

func (cfg *apiConfig) handleGetChirp(w http.ResponseWriter, r *http.Request) {
	pathVal := r.PathValue("id")
	id, err := uuid.Parse(pathVal)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id")
		return
	}

	dbChirp, err := cfg.db.GetChirpByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithJSON(w, http.StatusNotFound, map[string]string{})
			return
		}
		respondWithError(w, http.StatusInternalServerError, "fail to get chirp: "+err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, ChirpResponse{
		ID:        dbChirp.ID,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
	})
}
