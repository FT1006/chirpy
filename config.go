package main

import (
	"sync/atomic"

	"github.com/FT1006/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	Platform       string
	jwtSecret      string
}
