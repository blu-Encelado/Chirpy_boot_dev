package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	serv_mux := http.NewServeMux()
	apiCfg := apiConfig{}

	serv_mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	serv_mux.HandleFunc("GET /healthz", handlerReadiness)
	serv_mux.HandleFunc("GET /metrics", apiCfg.handleCountHits)
	serv_mux.HandleFunc("POST /reset", apiCfg.handleResetCountHits)

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

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})

}

func (cfg *apiConfig) handleCountHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	string_var := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	w.Write([]byte(string_var))
}

func (cfg *apiConfig) handleResetCountHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Swap(0)
	w.Write([]byte("Hits Resetted"))
}
