package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type RefreshToken struct {
	Token  string    `json:"token"`
	UserId uuid.UUID `json:"user_id"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	ExpiresAt string `json:"expires_at"`
	RevokedAt string `json:"revoked_at"`
}

func (cfg *apiConfig) handleListRefreshTokens(w http.ResponseWriter, r *http.Request) {
	refreshTokensDB, err := cfg.db.ListRefreshTokens(r.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(204)
			return
		}
		log.Printf("error to list refresh tokens %v", err)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	var response = []RefreshToken{}
	for _, tk := range refreshTokensDB {
		revokedDate := ""
		if !tk.RevokedAt.Time.IsZero() {
			revokedDate = tk.RevokedAt.Time.String()
		}
		response = append(response, RefreshToken{
			Token:     tk.Token,
			UserId:    tk.UserID,
			CreatedAt: tk.CreatedAt.String(),
			UpdatedAt: tk.UpdatedAt.String(),
			ExpiresAt: tk.ExpiresAt.String(),
			RevokedAt: revokedDate,
		})
	}

	respondWithJSON(w, http.StatusOK, response)
}
