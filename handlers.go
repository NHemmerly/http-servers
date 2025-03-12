package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/NHemmerly/http-servers/internal/auth"
	"github.com/NHemmerly/http-servers/internal/database"
	"github.com/google/uuid"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dB             *database.Queries
	platform       string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

type request struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	req := decodeRequest(w, r, &request{})
	user, err := cfg.getUserByEmail(req.Email, r)
	if err != nil {
		log.Printf("user not found")
		w.WriteHeader(401)
		return
	}
	err = auth.CheckPasswordHash(req.Password, user.HashedPassword)
	if err != nil {
		log.Printf("incorrect email or password")
		w.WriteHeader(401)
		return
	}
	responseWithJson(w, 200, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}

func (cfg *apiConfig) postChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	chirp, err := cfg.dB.CreateChirp(r.Context(), database.CreateChirpParams{
		UserID: params.UserId,
		Body:   params.Body,
	})
	if err != nil {
		log.Printf("could not create chirp: %s", err)
		return
	}
	responseWithJson(w, 201, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.dB.GetChirps(r.Context())
	if err != nil {
		log.Printf("could not retrieve all users: %s", err)
		return
	}
	var chirpArray []Chirp
	for _, chirp := range chirps {
		chirpArray = append(chirpArray, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		})
	}
	responseWithJson(w, 200, chirpArray)
}

func (cfg *apiConfig) getChirpById(w http.ResponseWriter, r *http.Request) {
	chirpId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		log.Printf("could not parse chirpId string: %s", err)
		return
	}
	chirp, err := cfg.dB.GetChirpByID(r.Context(), chirpId)
	if err != nil {
		log.Printf("chirp not found: %s", err)
		respondWithError(w, 404, "chirp not found")
	}
	responseWithJson(w, 200, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, req *http.Request) {
	params := decodeRequest(w, req, &request{})
	if params.Password == "" {
		log.Printf("No password provided")
		w.WriteHeader(400)
		return
	}
	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("could not hash password")
		w.WriteHeader(500)
		return
	}
	user, err := cfg.dB.CreateUser(req.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hash,
	})
	if err != nil {
		log.Printf("could not create user: %s", err)
		return
	}
	responseWithJson(w, 201, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email})

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
