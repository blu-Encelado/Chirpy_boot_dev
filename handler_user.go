package main

import (
	"Chirpy/internal/auth"
	"Chirpy/internal/database"
	"context"
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerRegisterUser(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	ctx := context.Background()

	hashed_pass, err := auth.HashPassword(params.Password)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	created_user := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashed_pass,
	}

	user_struct, err := cfg.db.CreateUser(ctx, created_user)

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

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hash, err := cfg.db.CheckPassword(ctx, params.Email)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't check password", err)
		return
	}

	is_valid, err := auth.CheckPasswordHash(params.Password, hash)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "fail to check", err)
		return
	}

	if is_valid == false {
		respondWithError(w, http.StatusUnauthorized, "not valid password", err)
		return
	}

	db_user, err := cfg.db.GetUser(ctx, params.Email)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't find user", err)
		return
	}

	user := User{
		ID:        db_user.ID,
		CreatedAt: db_user.CreatedAt,
		UpdatedAt: db_user.UpdatedAt,
		Email:     db_user.Email,
	}

	respondWithJson(w, http.StatusOK, user)
}
