package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
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

	cleanedBody := cleanBody(params.Body)
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}
	respondWithJSON(w, http.StatusOK, returnVals{CleanedBody: cleanedBody})

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
