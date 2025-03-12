package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

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
