package main

import (
	"Chirpy/internal/auth"
	"Chirpy/internal/database"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
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
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Fail to get token", nil)
		return
	}
	user_id, err := auth.ValidateJWT(token, secret_key)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", nil)
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
		UserID: user_id,
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
	request := []returnValues{}
	var chirps []database.Chirp
	var err error

	author_id := r.URL.Query().Get("author_id")
	if author_id != "" {
		author_uuid, err := uuid.Parse(author_id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Fail to get author_id", err)
			return
		}
		chirps, err = cfg.db.GetChirpsByAuthor(ctx, author_uuid)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Fail to get chirps by author_id", err)
			return
		}

	} else {
		chirps, err = cfg.db.GetAllChirps(ctx)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Fail to get chirps", err)
			return
		}
	}

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
	order := r.URL.Query().Get("sort")

	if order == "desc" {
		sort.Slice(request, func(i, j int) bool {
			return request[i].CreatedAt.After(request[j].CreatedAt)
		})
	} else {
		sort.Slice(request, func(i, j int) bool {
			return request[i].CreatedAt.Before(request[j].CreatedAt)
		})
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

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "coudn't find the token", err)
		return
	}
	user_id, err := auth.ValidateJWT(token, secret_key)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid jwt token", err)
		return
	}
	ctx := context.Background()
	id_requested := r.PathValue("chirpID")
	chirp_uuid, err := uuid.Parse(id_requested)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Fail to get chirp id", err)
		return
	}
	chirp, err := cfg.db.GetChirpFromId(ctx, chirp_uuid)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found", err)
		return
	}

	if user_id != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "the user is not the owner", errors.New("Unauthorized"))
		return
	}
	err = cfg.db.DeleteSingleChirp(ctx, chirp_uuid)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error on delete chirp", err)
		return
	}
	respondWithJson(w, http.StatusNoContent, nil)
}
