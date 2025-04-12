package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/FT1006/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error getting token")
		return
	}

	fmt.Println("Extracted Refresh Token:", refreshToken)

	dbRefreshToken, err := cfg.dbQueries.GetRefreshTokenByToken(r.Context(), refreshToken)
	if err != nil || dbRefreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Error getting refresh token")
		return
	}

	if dbRefreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token has been revoked")
		return
	}

	token, err := auth.MakeJWT(dbRefreshToken.UserID, cfg.jwtSecret, 3600*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		Token: token,
	})
}
