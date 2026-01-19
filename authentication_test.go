package main

import (
	"Chirpy/internal/auth"
	"testing"
)

func TestPassword(t *testing.T) {
	pass_1 := "12345"
	_, err := auth.HashPassword(pass_1)
	if err != nil {
		t.Errorf("fail to hash test: %v", err)
	}
}
