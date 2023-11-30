package main

import (
	"net/http"
)

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

func main() {
	//create new server multiplexer
	mux := http.NewServeMux()
	//run it through middleware for Cors header change
	corsMux := middlewareCors(mux)
	// assign it as http handler for root
	http.Handle("/", corsMux)
	//port to work on
	http.ListenAndServe(":8080", nil)
}
