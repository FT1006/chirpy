package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FT1006/chirpy/internal/auth"
	"github.com/FT1006/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error getting token")
		return
	}

	fmt.Println("Extracted Token:", token)
	fmt.Println("Loaded jwtSecret:", cfg.jwtSecret)
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		fmt.Println("Validation Error:", err)
		respondWithError(w, http.StatusUnauthorized, "Error validating token")
		return
	}

	fmt.Println("Validated UserID:", userID)

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}

	hashedPW, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error hashing password: %s", err))
		return
	}

	if err := cfg.dbQueries.UpdateUserEmailAndPassword(r.Context(), database.UpdateUserEmailAndPasswordParams{
		Email:          params.Email,
		HashedPassword: hashedPW,
		ID:             userID,
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error updating user: %s", err))
		return
	}

	dbUpdatedUser, err := cfg.dbQueries.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting user: %s", err))
		return
	}

	updatedUser := User{
		ID:          dbUpdatedUser.ID,
		CreatedAt:   dbUpdatedUser.CreatedAt,
		UpdatedAt:   dbUpdatedUser.UpdatedAt,
		Email:       dbUpdatedUser.Email,
		IsChirpyRed: dbUpdatedUser.IsChirpyRedTeam,
	}

	respondWithJSON(w, 200, updatedUser)
}
