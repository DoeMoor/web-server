package api

import (
	"fmt"
	"net/http"
)

func (cfg *ApiConfig) MetricsReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	cfg.FileserverHits.Store(0)
	fmt.Println("Middleware remove", cfg.FileserverHits.Load())
	w.Write([]byte("Hits reset to 0"))
}
