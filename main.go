package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/NHemmerly/http-servers/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("could not open database: %s", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	apiCfg := apiConfig{dB: dbQueries, platform: os.Getenv("PLATFORM")}
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	mux.HandleFunc("GET /api/healthz", func(writer http.ResponseWriter, req *http.Request) {
		req.Header.Set("Content-Type", "text/plain; charset=utf-8")
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("OK"))
	})

	mux.HandleFunc("GET /admin/metrics", apiCfg.getMetricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetMetricsHandler)
	mux.HandleFunc("POST /api/users", apiCfg.createUser)
	mux.HandleFunc("POST /api/chirps", apiCfg.postChirps)
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirps)
	mux.HandleFunc("GET /api/chirps/{id}", apiCfg.getChirpById)

	server.ListenAndServe()
}
