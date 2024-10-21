package auth_tests

import (
	"net/http"
	"testing"
	"time"

	"github.com/Kazyel/chirpy-bootdev/internal/auth"
	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userID, _ := uuid.NewRandom()

	token, err := auth.MakeJWT(userID, "secret", time.Hour)

	if err != nil {
		t.Errorf("Error making JWT: %s", err)
	}

	if len(token) != 208 {
		t.Errorf("Token length is not 256: %d", len(token))
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	token, err := auth.MakeJWT(userID, "secret", time.Hour)

	if err != nil {
		t.Errorf("Error making JWT: %s", err)
	}

	userIDFromToken, err := auth.ValidateJWT(token, "secret")

	if err != nil {
		t.Errorf("Error validating JWT: %s", err)
	}

	if userIDFromToken != userID {
		t.Errorf("User IDs do not match: %s != %s", userIDFromToken, userID)
	}
}

func TestValidateJWTWithInvalidToken(t *testing.T) {
	userID := uuid.New()

	token, _ := auth.MakeJWT(userID, "secret", time.Hour)

	_, err := auth.ValidateJWT(token, "not-secret")

	if err == nil {
		t.Errorf("Should have returned an error: Invalid Token")
	}
}

func TestVerboseJWTExpire(t *testing.T) {
	userID := uuid.New()
	token, err := auth.MakeJWT(userID, "secret", time.Second)

	if err != nil {
		t.Errorf("Error making JWT: %s", err)
	}

	time.Sleep(time.Second * 2)

	_, err = auth.ValidateJWT(token, "secret")
	t.Logf(err.Error())

	if err == nil {
		t.Errorf("Should have returned an error: Token expired")
	}
}

func TestVerboseBearerToken(t *testing.T) {
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+"token")

	token, err := auth.GetBearerToken(headers)
	t.Logf(token)

	if err != nil {
		t.Errorf("Error getting bearer token: %s", err)
	}

	if token != "token" {
		t.Errorf("Token does not match: %s != %s", token, "token")
	}
}
