package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"Chirpy/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	serv_mux := http.NewServeMux()
	apiCfg := apiConfig{}

	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	apiCfg.platform = os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbUrl)

	if err != nil {
		log.Printf("fail to opne .env: %s", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	apiCfg.db = dbQueries

	serv_mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	serv_mux.HandleFunc("GET /api/healthz", handlerReadiness)
	serv_mux.HandleFunc("GET /admin/metrics", apiCfg.handleCountHits)
	serv_mux.HandleFunc("POST /admin/reset", apiCfg.handleResetCountHits)
	serv_mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirps)

	serv_mux.HandleFunc("POST /api/users", apiCfg.handlerRegisterUser)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: serv_mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})

}

func (cfg *apiConfig) handleCountHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	string_var := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())
	w.Write([]byte(string_var))
}

func (cfg *apiConfig) handleResetCountHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Swap(0)
	w.Write([]byte("Metrics Resetted"))

	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "not dev", nil)
	}

	ctx := context.Background()
	err := cfg.db.ResetUser(ctx)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't reset users", err)
	}
}
