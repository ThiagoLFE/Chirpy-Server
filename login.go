package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github/ThiagoLFE/Chirpy-Server/internal/auth"
	"github/ThiagoLFE/Chirpy-Server/internal/database"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginResponse struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	email := strings.TrimSpace(loginRequest.Email)
	if len(email) == 0 {
		respondWithError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
			return
		}
		log.Printf("get user by email: %v", err)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	ok, err := auth.CheckPasswordHash(loginRequest.Password, user.HashedPassword)
	if err != nil {
		log.Printf("check password hash: \nuserEmail: %s\n err: %v", user.Email, err)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	refresh_token := auth.MakeRefreshToken()
	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refresh_token,
		UserID:    user.ID,
		ExpiresAt: time.Now().AddDate(0, 0, 60), //token lasts 60 days
	})

	if err != nil {
		log.Printf("Fail to create refresh token on db: %v", err)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.tokenSecret)
	if err != nil {
		log.Printf("fail to make JWT token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondWithJSON(w, http.StatusOK, LoginResponse{
		ID:           user.ID,
		Email:        user.Email,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		IsChirpyRed:  user.IsChirpyRed,
		Token:        token,
		RefreshToken: refresh_token,
	})
}
