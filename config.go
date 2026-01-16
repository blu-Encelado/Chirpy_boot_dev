package main

import (
	"Chirpy/internal/database"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	database       *database.Queries
}
