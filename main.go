package main

import (
	"database/sql"
	"github/ThiagoLFE/Chirpy-Server/internal/database"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	platform       string
	db             *database.Queries
	tokenSecret    string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	tokenSecret := os.Getenv("TOKEN_SECRET")

	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	dbQueries := database.New(dbConn)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		tokenSecret:    tokenSecret,
	}

	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsIncrement(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	mux.HandleFunc("POST /api/login", apiCfg.handleLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handleRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.handleRevokeRefreshToken)
	mux.HandleFunc("GET /api/refresh_tokens", apiCfg.handleListRefreshTokens)

	mux.HandleFunc("POST /api/users", apiCfg.handleCreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.MiddlewareAuth(apiCfg.handleUpdateUser))
	mux.HandleFunc("GET /api/users", apiCfg.handleListUsers)

	mux.HandleFunc("POST /api/chirps", apiCfg.MiddlewareAuth(apiCfg.handleCreateChirp))
	mux.HandleFunc("GET /api/chirps", apiCfg.handleListChirps)
	mux.HandleFunc("GET /api/chirps/{id}", apiCfg.handleGetChirp)
	mux.HandleFunc("DELETE /api/chirps/{id}", apiCfg.MiddlewareAuth(apiCfg.handleDeleteChirp))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
