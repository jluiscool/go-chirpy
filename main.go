package main

import (
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"
	//create new server multiplexer
	mux := http.NewServeMux()
	//make mux handle request to root path
	//allows the fileserver to serve index.html without returning index.html
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot))))
	mux.HandleFunc("/healthz", readynessEndPoint)
	//run it through middleware for Cors header change
	corsMux := middlewareCors(mux)
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

func readynessEndPoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	body := "OK"
	byteSlice := []byte(body)
	w.Write(byteSlice)
}
