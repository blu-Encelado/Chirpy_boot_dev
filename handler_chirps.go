package main

import (
	"Chirpy/internal/database"
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body    string    `json:"body"`
		User_id uuid.UUID `json:"user_id"`
	}
	type returnValues struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		User_id   uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp too long", nil)
		return
	}

	clean_string := cleanBadWords(params.Body)

	ctx := context.Background()
	chirp, err := cfg.db.CreateChirp(ctx, database.CreateChirpParams{
		Body:   clean_string,
		UserID: params.User_id,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Fail to create chirp", nil)
		return
	}
	request := returnValues{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		User_id:   chirp.UserID,
	}

	respondWithJson(w, http.StatusCreated, request)
}

func cleanBadWords(sentence string) string {
	bad_words := []string{"kerfuffle", "sharbert", "fornax"}
	list_words := strings.Split(sentence, " ")
	cleaned_words := []string{}

	for _, word := range list_words {
		low_word := strings.ToLower(word)
		if slices.Contains(bad_words, low_word) {
			cleaned_words = append(cleaned_words, "****")
		} else {
			cleaned_words = append(cleaned_words, word)
		}
	}
	return strings.Join(cleaned_words, " ")
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	type returnValues struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		User_id   uuid.UUID `json:"user_id"`
	}

	ctx := context.Background()
	chirps, err := cfg.db.GetAllChirps(ctx)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Fail to get chirps", err)
		return
	}

	request := []returnValues{}

	for _, chirp := range chirps {
		item := returnValues{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			User_id:   chirp.UserID,
		}
		request = append(request, item)
	}
	respondWithJson(w, http.StatusOK, request)
}

func (cfg *apiConfig) handlerGetSigleChirp(w http.ResponseWriter, r *http.Request) {
	type returnValues struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		User_id   uuid.UUID `json:"user_id"`
	}
	ctx := context.Background()
	id_requested := r.PathValue("chirpID")
	uuidRequested, err := uuid.Parse(id_requested)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Fail to get chirp id", err)
		return
	}
	chirp, err := cfg.db.GetChirpFromId(ctx, uuidRequested)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Fail to get chirp from id", err)
		return
	}
	request := returnValues{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		User_id:   chirp.UserID,
	}
	respondWithJson(w, http.StatusOK, request)
}
