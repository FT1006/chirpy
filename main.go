package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(http.StatusText(http.StatusOK)))
	if err != nil {
		log.Fatal(err)
	}
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hitCount := cfg.fileserverHits.Load()
	log.Printf("Hits: %d", hitCount)

	_, err := w.Write([]byte(fmt.Sprintf("Hits: %d", hitCount)))
	if err != nil {
		log.Fatal(err)
	}
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	log.Printf("Reset hits")

	_, err := w.Write([]byte("Hits reset to 0"))
	if err != nil {
		log.Fatal(err)
	}
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
	
func main() {
	filepathRoot := "."
	port := "8080"

	apiCfg := apiConfig{}

	serMux := http.NewServeMux()
	serMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	server := http.Server{
		Addr:    ":" + port,
		Handler: serMux,
	}

	serMux.HandleFunc("GET /healthz", handlerReadiness)
	serMux.HandleFunc("GET /metrics", apiCfg.handlerMetrics)
	serMux.HandleFunc("POST /reset", apiCfg.handlerReset)

	log.Printf("Serving files from %s on port %s", filepathRoot, port)
	log.Fatal(server.ListenAndServe())

}