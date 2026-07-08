package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	cfg := &apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux.HandleFunc("GET /api/healthz", readinessServe)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerNumberRequests)
	mux.HandleFunc("POST /admin/reset", cfg.handlerResetCountRequests)
	mux.HandleFunc("POST /api/validate_chirp", cfg.handlerValidateChirp)
	mux.Handle("/app/", cfg.middlewareCountRequest(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	log.Printf("Server is running at http://localhost%s", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func (c *apiConfig) middlewareCountRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (c *apiConfig) handlerNumberRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write(fmt.Appendf([]byte{}, `
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>
		`,
		c.fileserverHits.Load(),
	))
}

func readinessServe(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(200)

	w.Write([]byte("OK"))
}

func (c *apiConfig) handlerResetCountRequests(w http.ResponseWriter, _ *http.Request) {
	c.fileserverHits.Store(0)
	w.Write(fmt.Appendf([]byte{}, "reseted hits"))

}

func (c *apiConfig) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	var payload ValidateRequest

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
	}
	if len(payload.Body) <= 140 {
		respondWithJSON(w, http.StatusOK, ValidateResponse{
			Valid: true,
		})
		return
	}
	respondWithError(w, http.StatusBadRequest, "invalid body: max length must be at most of 140 characteres.")
}
