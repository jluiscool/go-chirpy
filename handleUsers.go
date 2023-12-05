package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jluiscool/go-chirpy/internal/database"
)

func handlerGetUsers(w http.ResponseWriter, r *http.Request) {
	//connection to db
	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Printf("error loading database: %s", err)
		w.WriteHeader(503)
		return
	}
	//get users func
	users, err := db.GetUsers()
	if err != nil {
		log.Printf("error getting users: %s", err)
		w.WriteHeader(500)
		return
	}
	//write headers
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	//handle json
	dat, err := json.Marshal(users)
	if err != nil {
		log.Printf("error encoding json: %s", err)
		w.WriteHeader(500)
	}
	//write json to response
	w.Write(dat)
}
