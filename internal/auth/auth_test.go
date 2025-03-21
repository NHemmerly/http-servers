package auth

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCompareHashPassword(t *testing.T) {
	// Testing table
	var tests = []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{"hash works", "thisisatestpassword", "thisisatestpassword", true},
		{"hash seed different", "ThisIsATestPassword", "sike", false},
	}
	// The execution loop
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hash, err := HashPassword(test.hash)
			if err != nil {
				t.Errorf("could not hash password: %s", err)
			}
			err = CheckPasswordHash(test.password, hash)
			if (err != nil) == test.want {
				t.Errorf("hash of %s does not match %s hash", test.password, hash)
			}
		})
	}

}

func TestValidateJWT(t *testing.T) {
	var tests = []struct {
		name      string
		duration  string
		secretKey string
		fakeKey   string
		uuid      string
		want      bool
	}{
		{"validation works", "10m", "thistestworks", "thistestworks", "f1e7d154-05fe-4dae-babd-e805734fe71b", true},
		{"expired key", "1ms", "thistestexpires", "thistestexpires", "f1e7d154-05fe-4dae-babd-e805734fe71b", false},
		{"wrong secret key", "10m", "correctKey", "wrongKey", "f1e7d154-05fe-4dae-babd-e805734fe71b", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			duration, err := time.ParseDuration(test.duration)
			if err != nil {
				t.Errorf("duration not parsed")
			}
			id, err := uuid.Parse(test.uuid)
			if err != nil {
				t.Errorf("uuid not parsed")
			}
			jwtStr, err := MakeJWT(id, test.secretKey, duration)
			if err != nil {
				t.Errorf("jwtString not made")
			}
			rightID, err := ValidateJWT(jwtStr, test.fakeKey)
			if (err != nil) == test.want {
				t.Errorf("could not validate jwt: %s", err)
			}
			if (rightID == id) != test.want {
				t.Errorf("RightID %v does not match original id %v", rightID, id)
			}
		})
	}

}

func TestGetBearerToken(t *testing.T) {
	var tests = []struct {
		name        string
		headerKey   string
		headerValue string
		want        bool
	}{
		{"header works", "Authorization", "Bearer tokenstring", true},
		{"header not there", "Authorization", "", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			header := make(http.Header)
			header.Add(test.headerKey, test.headerValue)
			tokenString := strings.Split(test.headerValue, " ")
			testVal, err := GetBearerToken(header)
			if (err != nil) == test.want {
				t.Errorf("could not get bearer token")
			}
			if len(tokenString) > 1 {
				if (testVal != tokenString[1]) == test.want {
					t.Errorf("%v does not equal %v to equal %v", testVal, tokenString[1], test.want)
				}
			}
		})
	}
}
