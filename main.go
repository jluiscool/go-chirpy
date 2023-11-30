package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	//create new server multiplexer
	mux := http.NewServeMux()
	//run it through middleware for Cors header change
	corsMux := middlewareCors(mux)

	//create server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())

	// // assign it as http handler for root
	// http.Handle("/", corsMux)
	// //port to work on
	// http.ListenAndServe(":8080", nil)
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
