package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github/ThiagoLFE/Chirpy-Server/internal/auth"
	"log"
	"net/http"
	"strings"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

	respondWithJSON(w, http.StatusOK, UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}
