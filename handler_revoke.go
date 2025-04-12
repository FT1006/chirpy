package main

import (
	"fmt"
	"net/http"

	"github.com/FT1006/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error getting token")
		return
	}

	if _, err := cfg.dbQueries.GetRefreshTokenByToken(r.Context(), refreshToken); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error getting refresh token")
		return
	}

	fmt.Println("Extracted Refresh Token:", refreshToken)

	if err := cfg.dbQueries.RevokeRefreshToken(r.Context(), refreshToken); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error revoking refresh token")
		return
	}

	w.WriteHeader(http.StatusNoContent) // This is 204
}
