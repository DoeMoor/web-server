package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/doemoor/wed-server/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) CreateChirp(w http.ResponseWriter, r *http.Request) {
	type chirpRequest struct {
		Body   string `json:"body"`
		UserId string `json:"user_id"`
	}

	type chirpResponse struct {
		Id        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    string    `json:"user_id"`
	}

	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	chirpReq := chirpRequest{}
	err := decoder.Decode(&chirpReq)
	if err != nil {
		err := responseWithError(w, 500, "Invalid JSON")
		if err != nil {
			log.Printf("error marshal error response :%v", err)
			return
		}
		return
	}

	if len(chirpReq.Body) > 140 {
		err := responseWithError(w, 400, "Chirp too long")
		if err != nil {
			log.Printf("response error :%v", err)
			return
		}
		return
	}

	cleanBody := clearBody(chirpReq.Body)
	userUUID, err := uuid.Parse(chirpReq.UserId)
	if err != nil {
		err := responseWithError(w, 400, "Invalid user id")
		if err != nil {
			log.Printf("response error :%v", err)
			return
		}
		return
	}

	reqParams := database.CreateChirpParams{
		Body: cleanBody,
		UserID: userUUID,
	}

	chirpDb, err := cfg.DbQueries.CreateChirp(r.Context(), reqParams)
	if err != nil {
		err := responseWithError(w, 500, "Error creating chirp")
		if err != nil {
			log.Printf("response error :%v", err)
			return
		}
		return
	}
	err = responseWithJson(w, 201, chirpResponse{
		Id:        chirpDb.ID.String(),
		CreatedAt: chirpDb.CreatedAt,
		UpdatedAt: chirpDb.UpdatedAt,
		Body:      chirpDb.Body,
		UserId:    chirpDb.UserID.String(),
	})
	if err != nil {
		log.Printf("error encoding response :%v", err)
		return
	}
}

var restrictedWords []string = []string{
	"kerfuffle",
	"sharbert",
	"fornax",
}

func clearBody(body string) string {
	splitBody := strings.Split(body, " ")
	for bodyIndex, bodyWord := range splitBody {
		for _, restrictedWord := range restrictedWords {
			if !strings.Contains(strings.ToLower(bodyWord), restrictedWord) {
				continue
			}
			splitBody[bodyIndex] = "****"
		}
	}
	return strings.Join(splitBody, " ")
}
