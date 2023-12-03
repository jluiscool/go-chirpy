package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const filepathRoot = "."
	// const adminPathRoot = "./admin"
	const port = "8080"
	apiCfg := apiConfig{
		fileserverHits: 0,
	}
	// adminCfg := apiConfig{
	// 	fileserverHits: 0,
	// }
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
	apiRouter.Post("/validate_chirp", postChirpValidation)
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
	log.Printf("Serving on port: %s\n", port)
	//ListenAndServe listens to TCP server.Addr, then calls Serve to handle incoming requests
	//main function blocks until the server is shut down, returning an error
	log.Fatal(srv.ListenAndServe())
}

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
	type returnValid struct {
		Valid bool `json:"valid"`
	}
	respBody := returnValid{
		Valid: true,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Write(dat)
}

func (cfg *apiConfig) getAdminIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf("<html>\n\n<body>\n<h1>Welcome, Chirpy Admin</h1>\n<p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits)))
	// cfg.fileserverHits++
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "no-cache")
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
}
