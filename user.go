package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRequest struct {
	Email string `json:"email"`
}

func (c *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var userRequest UserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode the parameters")
		return
	}

	email := strings.TrimSpace(userRequest.Email)

	if len(email) == 0 {
		respondWithError(w, http.StatusBadRequest, "email is required")
		return
	}

	dbUser, err := c.db.CreateUser(r.Context(), email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "fail to create user: "+err.Error())
		return
	}

	data := UserResponse{
		ID:        dbUser.ID,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}

	respondWithJSON(w, http.StatusCreated, data)
}

func (cfg *apiConfig) handleListUsers(w http.ResponseWriter, r *http.Request) {
	dbUsers, err := cfg.db.ListUsers(r.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithJSON(w, http.StatusOK, map[string]string{})
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	formattedList := make([]UserResponse, 0)
	for _, dbUsers := range dbUsers {
		formattedList = append(formattedList, UserResponse{
			ID:        dbUsers.ID,
			Email:     dbUsers.Email,
			CreatedAt: dbUsers.CreatedAt,
			UpdatedAt: dbUsers.UpdatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, formattedList)
}
