package main

import (
	"Chirpy/internal/auth"
	"Chirpy/internal/database"
	"context"
	"database/sql"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "fail to get token from the header", err)
		return
	}
	ctx := context.Background()
	err = cfg.db.SetRevokeFromToken(ctx, database.SetRevokeFromTokenParams{
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: time.Now(),
		Token:     refresh_token,
	})
	if err != nil {
		respondWithError(w, http.StatusBadGateway, "fail to change 'revoke_at'", err)
		return
	}
	respondWithJson(w, http.StatusNoContent, nil)
}
