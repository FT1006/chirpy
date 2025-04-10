package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FT1006/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}

	if dbUser, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	} else if err := auth.CheckPasswordHash(dbUser.HashedPassword, params.Password); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	} else {
		respondWithJSON(w, http.StatusOK, User{
			ID:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email:     dbUser.Email,
		})
	}
}
