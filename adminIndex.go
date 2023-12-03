package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) getAdminIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf("<html>\n\n<body>\n<h1>Welcome, Chirpy Admin</h1>\n<p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits)))
}
