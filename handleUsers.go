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

func handlePostUsers(w http.ResponseWriter, r *http.Request) {
	//make empty user
	type parameters struct {
		Email              string `json:"email"`
		Password           string `json:"password"`
		Expires_In_Seconds int    `json:"expires_in_seconds"`
	}
	//decode json
	//make new decoder
	decoder := json.NewDecoder(r.Body)
	//instanciate params
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	//get db connection
	dbCon, err := database.NewDB("./database.json")
	if err != nil {
		log.Printf("database connection error: %s", err)
		w.WriteHeader(500)
		return
	}
	//create new user
	newUser, err := dbCon.CreateUser(params.Email, params.Password)
	if err != nil {
		log.Printf("unable to create new user: %s", err)
		w.WriteHeader(500)
		return
	}
	//encode created user json
	dat, err := json.Marshal(newUser)
	if err != nil {
		log.Printf("unable to encode json: %s", err)
		w.WriteHeader(500)
		return
	}
	//write to response
	w.WriteHeader(201)
	w.Write(dat)
}

func handlePutUsers(w http.ResponseWriter, r *http.Request) {

}
