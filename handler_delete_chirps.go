package main

import (
	"fmt"
	"net/http"

	"github.com/FT1006/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteAChirp(w http.ResponseWriter, r *http.Request) {
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
	chirp, err := cfg.dbQueries.GetAChirp(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Chirp not found: %s", err)) // This is 404
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error getting token") // This is 401
		return
	}

	fmt.Println("Extracted Token:", token)
	fmt.Println("Loaded jwtSecret:", cfg.jwtSecret)
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		fmt.Println("Validation Error:", err)
		respondWithError(w, http.StatusUnauthorized, "Error validating token") // This is 401
		return
	}

	if chirp.UserID.String() != userID.String() {
		respondWithError(w, http.StatusForbidden, "You are not authorized to delete this chirp") // This is 403
		return
	}

	fmt.Println("Validated UserID:", userID)

	if err := cfg.dbQueries.DeleteAChirp(r.Context(), chirpUUID); err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting chirp: %s", err)) // This is 500
		return
	}

	w.WriteHeader(http.StatusNoContent) // This is 204
}
