package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jluiscool/go-chirpy/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	secret         string
}

func main() {
	//load env
	godotenv.Load()
	//load jwt secret
	jwtSecret := os.Getenv("JWT_SECRET")
	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}
	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
		secret:         jwtSecret,
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
	apiRouter.Post("/chirps", handlerPostChirp)
	apiRouter.Get("/chirps", handlerGetChirps)
	apiRouter.Get("/chirps/{id}", handlerGetChirp)
	//users endpoint
	apiRouter.Get("/users", handlerGetUsers)
	apiRouter.Post("/users", handlePostUsers)
	apiRouter.Put("/users", handlePutUsers)
	//login endpoint
	apiRouter.Post("/login", handleLogin)
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
	log.Printf("Serving on port: %s\n", port)
	//ListenAndServe listens to TCP server.Addr, then calls Serve to handle incoming requests
	//main function blocks until the server is shut down, returning an error
	log.Fatal(srv.ListenAndServe())
	//debug mode enabled
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	fmt.Println(dbg)
}
