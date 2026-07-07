package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(200)

		w.Write([]byte("OK"))
	})

	log.Printf("Server is running at http://localhost%s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
