package main

import (
	"log"
	"net/http"
	"time"

	"github.com/NHemmerly/http-servers/internal/auth"
)

type token struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) revokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refresh, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("could not retrieve refresh token: %s", err)
		return
	}
	if err = cfg.dB.RevokeToken(r.Context(), refresh); err != nil {
		log.Printf("could not revoke refresh token: %s", err)
		return
	}
	responseWithJson(w, 204, nil)
}

func (cfg *apiConfig) postRefreshToken(w http.ResponseWriter, r *http.Request) {
	refresh, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("could not retrieve refresh token: %s", err)
		return
	}
	refreshToken, err := cfg.dB.GetRefreshToken(r.Context(), refresh)
	if err != nil {
		log.Printf("Refresh token does not exist: %s", err)
		respondWithError(w, 401, "Invalid Token")
		return
	}
	currentTime := time.Now()
	if currentTime.Sub(refreshToken.ExpiresAt) >= 0 || refreshToken.RevokedAt.Valid {
		log.Printf("Refresh token expired")
		respondWithError(w, 401, "Invalid Token")
		return
	}
	newAccess, err := auth.MakeJWT(refreshToken.UserID, cfg.secret, time.Hour)
	if err != nil {
		log.Printf("could not make new jwt: %s", err)
	}
	responseWithJson(w, 200, token{
		Token: newAccess,
	})

}
