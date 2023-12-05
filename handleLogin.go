package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jluiscool/go-chirpy/internal/database"
)

func handleLogin(w http.ResponseWriter, r *http.Request) {
	//connection to db
	dbCon, err := database.NewDB("./database.json")
	if err != nil {
		log.Printf("error with db connection: %s", err)
		return
	}
	//get params
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := parameters{}
	//decode json
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding params: %s", err)
		w.WriteHeader(500)
		return
	}
	//look for user credentials
	user, err := dbCon.FindUser(params.Email, params.Password)
	if err != nil {
		log.Printf("error finding user: %s", err)
		w.WriteHeader(401)
		return
	}
	//encode json
	dat, err := json.Marshal(user)
	if err != nil {
		log.Printf("error encoding json: %s", err)
		w.WriteHeader(500)
		return
	}
	//if good login,, 200, else 404
	w.WriteHeader(200)
	w.Write(dat)
}
