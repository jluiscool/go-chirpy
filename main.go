package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{
		fileserverHits: 0,
	}
	fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	fsMidHandler := apiCfg.middlewareMetricsInc(fsHandler)
	//create new server multiplexer
	// mux := http.NewServeMux()
	r := chi.NewRouter()
	apiRouter := chi.NewRouter()
	adminRouter := chi.NewRouter()
	//make mux handle request to root path
	//allows the fileserver to serve index.html without returning index.html
	r.Handle("/app", fsMidHandler)
	r.Handle("/app/*", fsMidHandler)
	//api routes mounted
	r.Mount("/api", apiRouter)
	apiRouter.Get("/healthz", readynessEndPoint)
	apiRouter.Get("/metrics", apiCfg.handlerMetrics)
	apiRouter.Get("/reset", apiCfg.resetEndPoint)
	apiRouter.Post("/chirps", postChirpValidation)
	//admin routes mounted
	r.Mount("/admin", adminRouter)
	adminRouter.Get("/metrics", apiCfg.getAdminIndex)
	//run it through middleware for Cors header change
	corsMux := middlewareCors(r)
	//create server from a struct that describes the server configuration
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}
	//create DB
	NewDB("./database.json")
	log.Printf("Serving on port: %s\n", port)
	//ListenAndServe listens to TCP server.Addr, then calls Serve to handle incoming requests
	//main function blocks until the server is shut down, returning an error
	log.Fatal(srv.ListenAndServe())
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type DB struct {
	Path string
	Mux  *sync.RWMutex
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	newDB := DB{
		Path: path,
		Mux:  &sync.RWMutex{},
	}
	errDB := os.WriteFile(path, []byte(""), 0666)
	if errDB != nil {
		return nil, errDB
	}
	return &newDB, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	// Makes filecontent into a io.Reader
	dat, err := os.ReadFile(db.Path)
	if err != nil {
		return []Chirp{}, errors.New("could not read the database")
	}
	//turn []bytes to io.Reader
	reader := bytes.NewReader(dat)
	//decode the JSON
	decoder := json.NewDecoder(reader)
	dbData := DBStructure{}
	decoderErr := decoder.Decode(&dbData)
	if decoderErr != nil {
		log.Printf("Error decoding parameters: %s", err)
		return []Chirp{}, decoderErr
	}
	for _, newChirp := range dbData.Chirps {
		fmt.Println(newChirp.Body)
	}
	return []Chirp{}, nil
}
