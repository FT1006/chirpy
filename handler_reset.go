package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)

	_, err := w.Write([]byte("Hits reset to 0"))
	if err != nil {
		log.Fatal(err)
	}

	err = cfg.dbQueries.DeleteUsers(r.Context())
	if err != nil {
		log.Fatal(err)
	}
}
