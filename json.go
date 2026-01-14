package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println("err")
	}
	if code > 499 {
		log.Printf("Responding with 5xx error: %s", msg)
	}
	type errorReponse struct {
		Error string `json:"error"`
	}
	respondWithJson(w, code, errorReponse{
		Error: msg,
	})
}

func respondWithJson(w http.ResponseWriter, code int, paylod interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(paylod)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}
