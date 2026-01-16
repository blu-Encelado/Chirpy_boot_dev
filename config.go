package main

import (
	"Chirpy_boot_dev/main.go/internal/database"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	database       *database.Queries
}
