package main

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) getMetricsHandler(w http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := fmt.Sprintf(`<html>
  <body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())
	w.Write([]byte(hits))
}

func (cfg *apiConfig) resetMetricsHandler(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		log.Printf("function not available in non-dev environment")
		req.Header.Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(403)
		return
	}
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := cfg.dB.RemoveUsers(req.Context()); err != nil {
		log.Printf("could not remove users: %s", err)
	}
	cfg.fileserverHits.Store(0)
}
