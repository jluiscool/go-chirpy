package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type IDGenerator struct {
	counter int
}

type Chirp struct {
	Body string `json:"body"`
	Id   int    `json:"id"`
}

func postChirpValidation(w http.ResponseWriter, r *http.Request) {
	newCounter := IDGenerator{
		counter: 1,
	}
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
	w.WriteHeader(201)
	words := strings.Fields(params.Body)
	for i, word := range words {
		if strings.EqualFold(word, "kerfuffle") || strings.EqualFold(word, "sharbert") || strings.EqualFold(word, "fornax") {
			words[i] = "****"
		}
	}
	filteredSentence := strings.Join(words, " ")
	respBody := Chirp{
		Id:   newCounter.counter,
		Body: filteredSentence,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Write(dat)
}
