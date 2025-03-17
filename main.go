package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/NHemmerly/http-servers/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dB             *database.Queries
	platform       string
	secret         string
	polka          string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("could not open database: %s", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	apiCfg := apiConfig{
		dB:       dbQueries,
		platform: os.Getenv("PLATFORM"),
		secret:   os.Getenv("SECRET"),
		polka:    os.Getenv("POLKA_KEY"),
	}
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

	mux.HandleFunc("POST /admin/reset", apiCfg.resetMetricsHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.getMetricsHandler)
	mux.HandleFunc("PUT /api/users", apiCfg.updateLogin)
	mux.HandleFunc("POST /api/users", apiCfg.createUser)
	mux.HandleFunc("POST /api/login", apiCfg.loginUser)
	mux.HandleFunc("POST /api/chirps", apiCfg.postChirps)
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirps)
	mux.HandleFunc("GET /api/chirps/{id}", apiCfg.getChirpById)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.deleteChirp)
	mux.HandleFunc("POST /api/refresh", apiCfg.postRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.revokeRefreshToken)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.upgradeChirpyRed)

	server.ListenAndServe()
}
