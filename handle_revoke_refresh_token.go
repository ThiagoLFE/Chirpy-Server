package main

import (
	"database/sql"
	"errors"
	"github/ThiagoLFE/Chirpy-Server/internal/auth"
	"log"
	"net/http"
)

func (cfg *apiConfig) handleRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil || refreshToken == "" {
		respondWithError(w, http.StatusBadRequest, "refresh token is required on the authorization headers")
		return
	}

	if _, err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusBadRequest, "invalid refresh token")
			return
		}

		log.Printf("error to revoke refresh token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(204)
}
