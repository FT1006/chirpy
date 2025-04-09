package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/FT1006/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}

	const maxBodyLen = 140
	if len(params.Body) > maxBodyLen {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	fmt.Printf("Decoded ID: %v", params.UserID)

	cleanedBody := cleanBody(params.Body)
	chirpParams := database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: params.UserID,
	}

	fmt.Printf("Database params: %v", chirpParams)

	createdChirp, err := cfg.dbQueries.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating chirp: %s", err))
		return
	}

	fmt.Println("HTTP 201 Created")
	respondWithJSON(w, 201, Chirp{
		ID:        createdChirp.ID,
		CreatedAt: createdChirp.CreatedAt,
		UpdatedAt: createdChirp.UpdatedAt,
		Body:      createdChirp.Body,
		UserID:    createdChirp.UserID,
	})

	// TODO: Add chirp to user's chirps
	// TODO: Add user to chirp's users

	// TODO: Send email to user
	// TODO: Send email to admins

	// TODO: Add chirp to admin's chirps
	// TODO: Add admin to chirp's admins

	// TODO: Send email to admin

}

func cleanBody(body string) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		if strings.ToLower(word) == "kerfuffle" || strings.ToLower(word) == "sharbert" || strings.ToLower(word) == "fornax" {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	fmt.Printf("Cleaned body: %s\n", cleaned)
	return cleaned
}
