package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.dbQueries.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Error getting chirps: %s", err)) // This is 404
		return
	}

	chirps := make([]Chirp, len(dbChirps))
	for i, chirp := range dbChirps {
		chirps[i] = Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	respondWithJSON(w, http.StatusOK, chirps) // This is 200
}

func (cfg *apiConfig) handlerGetAChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		respondWithError(w, http.StatusNotFound, "chirpID is required") // This is 404
		return
	}
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Invalid chirpID: %s", err)) // This is 404
		return
	}

	dbChirp, err := cfg.dbQueries.GetAChirp(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Error getting chirp: %s", err)) // This is 404
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, chirp) // This is 200
}
