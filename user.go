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

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var userRequest User
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode the parameters")
		return
	}

	email := strings.TrimSpace(userRequest.Email)

	if len(email) == 0 {
		respondWithError(w, http.StatusBadRequest, "email is required")
		return
	}

	newUser, err := c.db.CreateUser(r.Context(), email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "fail to create user: "+err.Error())
		return
	}

	data := User{
		ID:        newUser.ID,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}

	respondWithJSON(w, http.StatusCreated, data)
}

func (cfg *apiConfig) handleListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := cfg.db.ListUsers(r.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithJSON(w, http.StatusOK, map[string]string{})
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	formattedList := make([]User, 0)
	for _, user := range users {
		formattedList = append(formattedList, User{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, formattedList)
}
