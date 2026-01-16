package main

import (
	"context"
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerRegisterUser(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	ctx := context.Background()
	user_struct, err := cfg.db.CreateUser(ctx, params.Email)

	user := User{
		ID:        user_struct.ID,
		CreatedAt: user_struct.CreatedAt,
		UpdatedAt: user_struct.UpdatedAt,
		Email:     user_struct.Email,
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't register user", err)
	}

	respondWithJson(w, http.StatusCreated, user)
}
