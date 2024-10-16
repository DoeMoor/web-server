package api

import (
	"fmt"
	"net/http"
	"os"
)

func (cfg *ApiConfig) MetricsReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if os.Getenv("PLATFORM") != "dev" {
		responseWithError(w, 403, "Forbidden")
	}

	w.WriteHeader(http.StatusOK)
	cfg.FileserverHits.Store(0)
	fmt.Println("Middleware remove", cfg.FileserverHits.Load())
	w.Write([]byte("\"Hits\" reset to 0, \"user\" schema is empty"))
}
