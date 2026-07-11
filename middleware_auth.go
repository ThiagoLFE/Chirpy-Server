package main

import (
	"context"
	"github/ThiagoLFE/Chirpy-Server/internal/auth"
	"net/http"
)

type contextKey string

const userIDContextKey contextKey = "userID"

func (cfg *apiConfig) MiddlewareAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		userID, err := auth.ValidateJWT(token, cfg.tokenSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), userIDContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
