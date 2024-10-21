package auth_tests

import (
	"testing"

	"github.com/Kazyel/chirpy-bootdev/internal/auth"
)

func TestHashPassword(t *testing.T) {
	password := "password"
	hash, err := auth.HashPassword(password)

	if err != nil {
		t.Errorf("Error hashing password: %s", err)
	}

	if len(hash) != 60 {
		t.Errorf("Hash length is not 60: %d", len(hash))
	}
}

func TestCheckPassword(t *testing.T) {
	password := "password"
	hash, err := auth.HashPassword(password)

	if err != nil {
		t.Errorf("Error hashing password: %s", err)
	}

	err = auth.CheckPasswordHash(password, hash)

	if err != nil {
		t.Errorf("Error checking password: %s", err)
	}
}
