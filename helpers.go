package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/NHemmerly/http-servers/internal/database"
)

type upgrade struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

type login struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type parameters struct {
	Body string `json:"body"`
}

func decodeRequest(w http.ResponseWriter, req *http.Request, form interface{}) error {
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(form)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return err
	}
	return nil
}

func (l *login) decodeRequest(w http.ResponseWriter, req *http.Request) error {
	return decodeRequest(w, req, l)
}

func (p *parameters) decodeRequest(w http.ResponseWriter, req *http.Request) error {
	return decodeRequest(w, req, p)
}

func (u *upgrade) decodeRequest(w http.ResponseWriter, req *http.Request) error {
	return decodeRequest(w, req, u)
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

func (cfg *apiConfig) getUserByEmail(email string, r *http.Request) (*database.User, error) {
	if user, err := cfg.dB.GetUserByEmail(r.Context(), email); err != nil {
		return nil, fmt.Errorf("could not get user by email: %w", err)
	} else {
		return &user, nil
	}
}
