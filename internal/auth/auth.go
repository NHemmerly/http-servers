package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GetAPIKey(headers http.Header) (string, error) {
	if bearer := headers.Get("Authorization"); bearer == "" {
		return "", fmt.Errorf("authorization header not found")
	} else {
		tokenString := strings.Split(bearer, " ")
		if len(tokenString) < 1 {
			return "", fmt.Errorf("no api key found")
		}
		if strings.ToLower(tokenString[0]) != "apikey" {
			return "", fmt.Errorf("no api key found")
		}
		return tokenString[1], nil
	}
}

func MakeRefreshToken() string {
	seed := make([]byte, 32)
	rand.Read(seed)
	hexString := hex.EncodeToString(seed)
	return hexString
}

func HashPassword(pass string) (string, error) {
	pword, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
	if err != nil {
		return "", fmt.Errorf("could not hash pass: %w", err)
	}

	return string(pword[:]), nil
}

func CheckPasswordHash(password, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return fmt.Errorf("password and hash do not match: %w", err)
	}
	return nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   fmt.Sprintf("%v", userID),
	})
	jwt, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("could not sign token: %w", err)
	}
	return jwt, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("could not parse: %w", err)
	}
	uuidString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, fmt.Errorf("could not get subject: %w", err)
	}
	id, err := uuid.Parse(uuidString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("could not parse uuid: %w", err)
	}
	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	if bearer := headers.Get("Authorization"); bearer == "" {
		return "", fmt.Errorf("authorization header not found")
	} else {
		tokenString := strings.Split(bearer, " ")
		if len(tokenString) < 1 {
			return "", fmt.Errorf("no bearer token found")
		}
		return tokenString[1], nil
	}
}
