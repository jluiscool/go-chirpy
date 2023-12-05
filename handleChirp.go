package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/jluiscool/go-chirpy/internal/database"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func handlerPostChirp(w http.ResponseWriter, r *http.Request) {
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

	cleanedChirp, err := validateChirp(params.Body)
	if err != nil {
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
	//write to database
	dbCon, err := database.NewDB("./database.json")
	if err != nil {
		log.Printf("Error handling db connection: %s", err)
		return
	}
	newChirp, err := dbCon.CreateChirp(cleanedChirp)
	if err != nil {
		log.Printf("Error creating new chirp: %s", err)
		return
	}
	//encode response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	respBody := newChirp
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Write(dat)
}

func validateChirp(body string) (string, error) {
	if len(body) > 140 {
		return "", errors.New("chirp is too long")
	}
	words := strings.Fields(body)
	for i, word := range words {
		if strings.EqualFold(word, "kerfuffle") || strings.EqualFold(word, "sharbert") || strings.EqualFold(word, "fornax") {
			words[i] = "****"
		}
	}
	filteredSentence := strings.Join(words, " ")
	return filteredSentence, nil
}

func handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	dbCon, err := database.NewDB("./database.json")
	if err != nil {
		log.Printf("Error handling db connection: %s", err)
		return
	}
	allChirps, err := dbCon.GetChirps()
	if err != nil {
		log.Printf("Error getting chirps: %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	dat, err := json.Marshal(allChirps)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Write(dat)
}
