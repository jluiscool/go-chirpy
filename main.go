package main

import (
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
