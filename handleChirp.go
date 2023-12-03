package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func postChirpValidation(w http.ResponseWriter, r *http.Request) {
	//decode request
	type parameters struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	//encode response
	w.Header().Set("Content-Type", "application/json")
	// log.Printf(strconv.Itoa(len(params.Body)))
	if len(params.Body) > 140 {
		type returnErr struct {
			Err string `json:"error"`
		}
		errRes := returnErr{
			Err: "chirp is too long",
		}
		dat, err := json.Marshal(errRes)
		if err != nil {
			log.Printf("Error sending error JSON: %s", err)
		}
		w.WriteHeader(400)
		w.Write(dat)
		return
	}
	w.WriteHeader(200)
	words := strings.Fields(params.Body)
	for i, word := range words {
		if strings.EqualFold(word, "kerfuffle") || strings.EqualFold(word, "sharbert") || strings.EqualFold(word, "fornax") {
			words[i] = "****"
		}
	}
	filteredSentence := strings.Join(words, " ")
	type returnValid struct {
		Cleaned_body string `json:"cleaned_body"`
	}
	respBody := returnValid{
		Cleaned_body: filteredSentence,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Write(dat)
}
