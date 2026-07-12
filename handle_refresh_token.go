package main

import (
	"database/sql"
	"errors"
	"github/ThiagoLFE/Chirpy-Server/internal/auth"
	"log"
	"net/http"
	"time"
)

type ResponseRefreshToken struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "refresh token is required")
		return
	}

	refreshTokenDB, err := cfg.db.GetRefreshTokenByRefreshToken(r.Context(), refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, "refresh token invalid")
			return
		}
		log.Printf("fail to get refresh token from db: %v", err)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if refreshTokenDB.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "refresh token expired, loggin again")
		return
	}

	if !refreshTokenDB.RevokedAt.Time.IsZero() {
		respondWithError(w, http.StatusUnauthorized, "refresh token revoked")
		return
	}

	accessToken, err := auth.MakeJWT(refreshTokenDB.UserID, cfg.tokenSecret)
	if err != nil {
		log.Printf("fail to create new access token on refresh: %v", err)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	respondWithJSON(w, http.StatusOK, ResponseRefreshToken{
		Token: accessToken,
	})
}
