package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCompareHashPassword(t *testing.T) {
	password := "thisisatestpassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("could not hash password: %s", err)
	}
	if err := CheckPasswordHash(password, hash); err != nil {
		t.Errorf("hash of %s does not match %s hash", password, hash)
	}
}

func TestValidateJWT(t *testing.T) {
	duration, err := time.ParseDuration("10m")
	if err != nil {
		t.Errorf("duration not parsed")
	}
	id, err := uuid.Parse("f1e7d154-05fe-4dae-babd-e805734fe71b")
	if err != nil {
		t.Errorf("uuid not parsed")
	}
	jwtStr, err := MakeJWT(id, "SecretKey123", duration)
	if err != nil {
		t.Errorf("jwtString not made")
	}
	rightID, err := ValidateJWT(jwtStr, "SecretKey123")
	if err != nil {
		t.Errorf("could not validate jwt")
	}
	if rightID != id {
		t.Errorf("RightID %v does not match original id %v", rightID, id)
	}
}
