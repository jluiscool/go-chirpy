package main

import "net/http"

func readynessEndPoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	body := "OK"
	byteSlice := []byte(body)
	w.Write(byteSlice)
}
