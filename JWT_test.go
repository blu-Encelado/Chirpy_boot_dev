package main

import (
	"Chirpy/internal/auth"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT_Success(t *testing.T) {
	secret := "super-secret"
	userID := uuid.New()

	token, err := auth.MakeJWT(userID, secret, time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	gotID, err := auth.ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT returned error: %v", err)
	}

	if gotID != userID {
		t.Fatalf("expected userID %v, got %v", userID, gotID)
	}
}

func TestExpiredJWT(t *testing.T) {
	secret := "super-secret"
	userID := uuid.New()

	token, err := auth.MakeJWT(userID, secret, time.Duration(-10))
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	_, err = auth.ValidateJWT(token, secret)

	if err == nil {
		t.Fatal("ValidateJWT validated an expired token")
	}
}

func TestWrongJWT(t *testing.T) {
	secret := "super-secret-a"
	userID := uuid.New()

	token, err := auth.MakeJWT(userID, secret, time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	_, err = auth.ValidateJWT(token, "super-secret-b")
	if err == nil {
		t.Fatal("ValidateJWT validated token with a different 'tokenSecret'")
	}
}
