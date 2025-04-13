package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FT1006/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error getting API key") // This is 401
		return
	}

	fmt.Println("Extracted API Key:", apiKey)
	fmt.Println("Loaded apiKey:", cfg.polkaApiKey)
	if apiKey != cfg.polkaApiKey {
		respondWithError(w, http.StatusUnauthorized, "Error validating API key") // This is 401
		return
	}

	var parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&parameters); err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}

	if parameters.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent) // This is 204
		return
	}

	UserUUID, err := uuid.Parse(parameters.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error parsing user ID: %s", err))
		return
	}

	if _, err := cfg.dbQueries.GetUserByID(r.Context(), UserUUID); err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("User not found: %s", err))
		return
	}

	if err := cfg.dbQueries.SetUserChirpyRed(r.Context(), UserUUID); err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error setting user as chirpy red: %s", err))
		return
	}

	w.WriteHeader(http.StatusNoContent) // This is 204
}
