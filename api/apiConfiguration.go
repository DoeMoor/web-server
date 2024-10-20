package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/doemoor/web-server/internal/database"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	DbQueries      *database.Queries
	Secret         string
	ApiKey         string
}

func responseWithError(w http.ResponseWriter, code int, message string) error {
	return responseWithJson(w, code, map[string]string{"error": message})
}

func responseWithJson(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Println("responseWithJson: error marshaling json: " + err.Error())
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}
