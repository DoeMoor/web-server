package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) metricsReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	fmt.Println("Middleware remove", cfg.fileserverHits.Load())
	w.Write([]byte("Hits reset to 0"))
}
