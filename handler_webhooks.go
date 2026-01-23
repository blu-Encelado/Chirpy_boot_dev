package main

import (
	"Chirpy/internal/auth"
	"Chirpy/internal/database"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerWebHooks(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "fail to get apiKey", err)
		return
	}
	if apiKey != polka_key {
		respondWithError(w, http.StatusUnauthorized, "apiKey != polkaKey", errors.New("polka missmatch"))
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			User_id uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters for webHooks", err)
		return
	}
	if params.Event != "user.upgraded" {
		respondWithJson(w, http.StatusNoContent, nil)
		return
	}
	ctx := context.Background()
	err = cfg.db.UpdateIsChirpyRed(ctx, database.UpdateIsChirpyRedParams{
		IsChirpyRed: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
		ID: params.Data.User_id,
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Fail to update is_chirpy_red", err)
		return
	}
	respondWithJson(w, http.StatusNoContent, nil)

}
