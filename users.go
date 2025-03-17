package main

import (
	"log"
	"net/http"
	"time"

	"github.com/NHemmerly/http-servers/internal/auth"
	"github.com/NHemmerly/http-servers/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	var req login
	if err := req.decodeRequest(w, r); err != nil {
		respondWithError(w, http.StatusInternalServerError, "server error")
		return
	}
	user, err := cfg.getUserByEmail(req.Email, r)
	if err != nil {
		respondWithError(w, 401, "user not found")
		return
	}
	if err = auth.CheckPasswordHash(req.Password, user.HashedPassword); err != nil {
		respondWithError(w, 401, "incorrect email or password")
		return
	}
	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		log.Printf("could not create jwt token")
		respondWithError(w, http.StatusInternalServerError, "server error")
		return
	}
	newRefToken, err := cfg.dB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  auth.MakeRefreshToken(),
		UserID: user.ID,
	})
	if err != nil {
		log.Printf("could not create RefreshToken: %s", err)
		respondWithError(w, http.StatusInternalServerError, "server error")
		return
	}
	responseWithJson(w, 200, User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: newRefToken.Token,
		IsChirpyRed:  user.IsChirpyRed,
	})
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, req *http.Request) {
	params := login{}
	params.decodeRequest(w, req)
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
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed})

}

func (cfg *apiConfig) updateLogin(w http.ResponseWriter, r *http.Request) {
	access, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("could not retrieve access token: %s", err)
		respondWithError(w, http.StatusUnauthorized, "no access token")
		return
	}
	var creds login
	if err := creds.decodeRequest(w, r); err != nil {
		log.Printf("could not decode request: %s", err)
		respondWithError(w, http.StatusInternalServerError, "server error")
		return
	}
	uuid, err := auth.ValidateJWT(access, cfg.secret)
	if err != nil {
		log.Printf("could not validate user: %s", err)
		respondWithError(w, http.StatusUnauthorized, "unauthorized user")
		return
	}
	hashed, err := auth.HashPassword(creds.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not hash password")
		return
	}
	user, err := cfg.dB.UpdateCreds(r.Context(), database.UpdateCredsParams{
		HashedPassword: hashed,
		Email:          creds.Email,
		ID:             uuid,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not update db")
		return
	}
	responseWithJson(w, http.StatusOK, User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})

}

func (cfg *apiConfig) upgradeChirpyRed(w http.ResponseWriter, r *http.Request) {
	var upgrade upgrade
	if err := upgrade.decodeRequest(w, r); err != nil {
		respondWithError(w, http.StatusInternalServerError, "server error")
		return
	}
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not find api key")
		return
	}
	if apiKey != cfg.polka {
		respondWithError(w, http.StatusUnauthorized, "wrong api key")
		return
	}
	user_id, err := uuid.Parse(upgrade.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not parse path")
		return
	}
	if upgrade.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "ignored event")
		return
	}
	if err := cfg.dB.UpgradeUser(r.Context(), user_id); err != nil {
		respondWithError(w, http.StatusNotFound, "user not found")
		return
	}
	respondWithError(w, http.StatusNoContent, "user updated")
}
