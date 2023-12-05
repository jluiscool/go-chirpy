package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

type Chirp struct {
	Body string `json:"body"`
	Id   int    `json:"id"`
}

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

// NewDB creates new db map connection
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

// ensures db is read
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	//if file doesn't exist
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps: map[int]Chirp{},
		Users:  map[int]User{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	//to write to dbStructure, encode to JSON, which turns to []byte
	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	//write JSON to file, gets err if any
	err = os.WriteFile(db.path, dat, 0600)
	if err != nil {
		return err
	}
	return nil
}

// called if ensureDB sees its been read
func (db *DB) loadDB() (DBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStructure := DBStructure{}
	dat, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}
	//gets dat to decode the json and sets it to dbStructure, returns err if any
	err = json.Unmarshal(dat, &dbStructure)
	if err != nil {
		return dbStructure, err
	}
	return dbStructure, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirpByID(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	foundChirp := Chirp{}
	for _, chirp := range dbStructure.Chirps {
		if chirp.Id == id {
			foundChirp = chirp
		}
	}
	if foundChirp.Id == 0 {
		return foundChirp, errors.New("no chirp found")
	}
	return foundChirp, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	//load db
	dbStructure, err := db.loadDB()
	//return empty if err
	if err != nil {
		return Chirp{}, err
	}
	//get new id
	newId := len(dbStructure.Chirps) + 1
	//make chirp
	newChirp := Chirp{
		Body: body,
		Id:   newId,
	}
	//put new chirp to dbStructure
	dbStructure.Chirps[newId] = newChirp
	//write new structure to file
	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}
	return newChirp, nil
}

func (db *DB) GetUsers() ([]User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return []User{}, errors.New("error loading database")
	}
	users := make([]User, 0, len(dbStructure.Users))
	for _, user := range dbStructure.Users {
		users = append(users, user)
	}
	return users, nil
}

func (db *DB) CreateUser(email string) (User, error) {
	//load db
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	//get new id
	userId := len(dbStructure.Users) + 1
	//make user
	user := User{
		Email: email,
		Id:    userId,
	}
	//put user into structure
	dbStructure.Users[userId] = user
	//put structure into db
	err = db.writeDB(dbStructure)
	if err != nil {
		return user, err
	}
	//return user
	return user, nil
}
