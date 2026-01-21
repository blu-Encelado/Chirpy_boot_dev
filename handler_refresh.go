package main

import (
	"Chirpy/internal/auth"
	"context"
	"errors"
	"net/http"
	"time"
)

var ErrTokenExpired = errors.New("token expired")

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "fail to get token from the header", err)
		return
	}
	ctx := context.Background()
	token_value, err := cfg.db.GetValuesFromToken(ctx, refresh_token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "fail to get token from db", err)
		return
	}

	if time.Until(token_value.ExpiresAt) <= 0 {
		respondWithError(w, http.StatusUnauthorized, "refresh: token expired", ErrTokenExpired)
		return
	}
	if token_value.RevokedAt.Valid == true {
		respondWithError(w, http.StatusUnauthorized, "refresh: token revoked", ErrTokenExpired)
		return
	}
	type ref_token struct {
		Token string `json:"token"`
	}
	new_jwt, err := auth.MakeJWT(token_value.UserID, secret_key, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "fail to create jwt", err)
		return
	}
	last_token := ref_token{
		Token: new_jwt,
	}
	respondWithJson(w, http.StatusOK, last_token)
}
