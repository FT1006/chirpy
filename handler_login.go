package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/FT1006/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiredInSeconds int    `json:"expired_in_seconds"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}

	if params.ExpiredInSeconds == 0 || params.ExpiredInSeconds > 3600 {
		params.ExpiredInSeconds = 3600
	}

	if dbUser, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	} else if err := auth.CheckPasswordHash(dbUser.HashedPassword, params.Password); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	} else {
		token, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, time.Duration(params.ExpiredInSeconds)*time.Second)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error generating token")
			return
		}
		respondWithJSON(w, http.StatusOK, User{
			ID:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email:     dbUser.Email,
			Token:     token,
		})
	}
}
