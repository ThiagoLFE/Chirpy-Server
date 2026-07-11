package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github/ThiagoLFE/Chirpy-Server/internal/auth"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type LoginRequest struct {
	Email            string        `json:"email"`
	Password         string        `json:"password"`
	ExpiresInSeconds time.Duration `json:"expires_in_seconds"`
}
type LoginResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
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

	if loginRequest.ExpiresInSeconds.String() == "0s" {
		loginRequest.ExpiresInSeconds = time.Duration(3600) // seconds to be an hour
	}

	token, err := auth.MakeJWT(user.ID, cfg.tokenSecret, loginRequest.ExpiresInSeconds*time.Second)
	if err != nil {
		log.Printf("fail to make JWT token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondWithJSON(w, http.StatusOK, LoginResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Token:     token,
	})
}
