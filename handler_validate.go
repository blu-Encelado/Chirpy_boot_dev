package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func handlerChirpValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnValues struct {
		Cleaned_body string `json:"cleaned_body"`
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

	respondWithJson(w, http.StatusOK, returnValues{
		Cleaned_body: clean_string,
	})
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
