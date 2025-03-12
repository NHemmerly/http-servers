package auth

import (
	"fmt"

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
