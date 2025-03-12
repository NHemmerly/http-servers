package auth

import "testing"

func testCompareHashPassword(t *testing.T) {
	password := "thisisatestpassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("could not hash password: %s", err)
	}
	if err := CheckPasswordHash(password, hash); err != nil {
		t.Errorf("hash of %s does not match %s hash", password, hash)
	}
}
