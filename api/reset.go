package api

import (
	"log"
	"net/http"
	"os"
)

func (cfg *ApiConfig) MetricsReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if os.Getenv("PLATFORM") != "dev" {
		responseWithError(w, 403, "allow only on dev environment")
		return
	}

	if err := cfg.DbQueries.DeleteUsers(r.Context()); err != nil {
		log.Println("Error deleting users" + err.Error())
		responseWithError(w, 500, "Error deleting users")
		return
	}
	w.WriteHeader(http.StatusOK)
	cfg.FileserverHits.Store(0)
	w.Write([]byte("\"Hits\" reset to 0, \"user\" schema is empty"))
}
