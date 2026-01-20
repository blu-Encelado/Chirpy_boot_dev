package main

import (
	"Chirpy/internal/auth"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPassword(t *testing.T) {
	pass_1 := "12345"
	_, err := auth.HashPassword(pass_1)
	if err != nil {
		t.Errorf("fail to hash test: %v", err)
	}
}

func TestToken(t *testing.T) {
	header := http.Header{}
	user := uuid.New()
	s_key := "secret"
	token_raw, err := auth.MakeJWT(user, s_key, time.Duration(50))
	if err != nil {
		t.Fatalf("Fail to MakeJWT: %v", err)
	}
	header.Set("Authorization", "Bearer "+token_raw)

	token, err := auth.GetBearerToken(header)
	if err != nil {
		t.Fatalf("Fail to check token: %v", err)
	}
	if token == "" {
		t.Fatal("nil token")
	}
}
