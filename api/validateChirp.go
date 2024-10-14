package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func ValidateChirp(w http.ResponseWriter, r *http.Request) {
	type chirpRequest struct {
		Body string `json:"body"`
	}

	type chirpError struct {
		ChirpErr string `json:"error"`
	}

	type chirpResponse struct {
		Valid bool `json:"valid"`	
	}

	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	chirpReq := chirpRequest{}
	err := decoder.Decode(&chirpReq)
	if err != nil {
		data, err := json.Marshal(chirpError{ChirpErr: "Invalid JSON"})
		if err != nil {
			log.Printf("error marshal error response :%v", err)
			return
		}
		w.WriteHeader(500)
		w.Write(data)
		return
	}
	if len(chirpReq.Body) > 140 {
		data, err := json.Marshal(chirpError{ChirpErr: "Chirp is too long"})
		if err != nil {
			log.Printf("error marshal error response :%v", err)
			return
		}
		w.WriteHeader(400)
		w.Write(data)
		return
	}

	data, err := json.Marshal(chirpResponse{Valid: true })
	if err != nil {
		log.Printf("error marshal error response :%v", err)
		return
	}
	w.WriteHeader(200)
	w.Write(data)
}
