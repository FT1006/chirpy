package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/FT1006/chirpy/internal/auth"
	"github.com/FT1006/chirpy/internal/database"
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
		token, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, 3600*time.Second)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error generating token")
			return
		}
		refreshToken, err := auth.MakeRefreshToken()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error generating refresh token")
			return
		}

		dbRefreshToken, err := cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
			Token:  refreshToken,
			UserID: dbUser.ID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error creating refresh token")
			return
		}

		fmt.Println(dbRefreshToken)

		respondWithJSON(w, http.StatusOK, User{
			ID:           dbUser.ID,
			CreatedAt:    dbUser.CreatedAt,
			UpdatedAt:    dbUser.UpdatedAt,
			Email:        dbUser.Email,
			Token:        token,
			RefreshToken: refreshToken,
		})
	}
}
