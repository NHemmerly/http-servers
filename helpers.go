package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/NHemmerly/http-servers/internal/database"
)

func parseExpiration(r *request) {
	if r.ExpiresInSeconds == 0 {
		r.ExpiresInSeconds = 3600
	} else if r.ExpiresInSeconds > 3600 {
		r.ExpiresInSeconds = 3600
	}
}

func decodeRequest(w http.ResponseWriter, req *http.Request, params *request) *request {
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return nil
	}
	return params
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type returnVals struct {
		Error string `json:"error"`
	}
	respBody := returnVals{
		Error: msg,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func responseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func cleanDirty(w http.ResponseWriter, code int, msg string) {
	words := strings.Split(msg, " ")
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	for i, word := range words {
		if slices.Contains(badWords, strings.ToLower(word)) {
			words[i] = "****"
		}
	}
	resp := strings.Join(words, " ")
	responseWithJson(w, code, map[string]string{
		"cleaned_body": resp,
	})
}

func (cfg *apiConfig) getUserByEmail(email string, r *http.Request) (*database.User, error) {
	if user, err := cfg.dB.GetUserByEmail(r.Context(), email); err != nil {
		return nil, fmt.Errorf("could not get user by email: %w", err)
	} else {
		return &user, nil
	}
}
