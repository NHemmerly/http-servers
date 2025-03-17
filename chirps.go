package main

import (
	"log"
	"net/http"
	"time"

	"github.com/NHemmerly/http-servers/internal/auth"
	"github.com/NHemmerly/http-servers/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
	access, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no access token")
		return
	}
	user_id, err := auth.ValidateJWT(access, cfg.secret)
	if err != nil {
		log.Printf("could not validate user: %s", err)
		respondWithError(w, http.StatusUnauthorized, "unauthorized user")
		return
	}
	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not parse path")
		return
	}
	chirp, err := cfg.dB.GetChirpByID(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found")
		return
	}
	if chirp.UserID != user_id {
		respondWithError(w, http.StatusForbidden, "Unauthorized")
		return
	}
	if err := cfg.dB.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     chirpId,
		UserID: user_id,
	}); err != nil {
		respondWithError(w, http.StatusForbidden, "Unauthorized")
		return
	}
	responseWithJson(w, http.StatusNoContent, "Chirp deleted")
}

func (cfg *apiConfig) postChirps(w http.ResponseWriter, r *http.Request) {
	params := parameters{}
	params.decodeRequest(w, r)
	// bearer and token auth
	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("could not get bearer: %s", err)
		return
	}
	uuid, err := auth.ValidateJWT(bearer, cfg.secret)
	if err != nil {
		log.Printf("%s", uuid)
		log.Printf("could not validate jwt: %s", err)
		respondWithError(w, 401, "Unauthorized")
		return
	}
	// chirp length verification
	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	chirp, err := cfg.dB.CreateChirp(r.Context(), database.CreateChirpParams{
		UserID: uuid,
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
		respondWithError(w, 404, "chirp not found")
		return
	}
	responseWithJson(w, 200, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
}
