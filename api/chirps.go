package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
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

func (cfg *ApiConfig) GetChirps(w http.ResponseWriter, r *http.Request) {
	type chirpResponse struct {
		Id        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    string    `json:"user_id"`
	}


	chirpsArray, err := cfg.DbQueries.GetAllChirps(r.Context())
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Printf("error fetching chirps from db :%v", err)
		responseWithError(w, 500, "Error fetching chirps")
		return
	}
	var chirps []chirpResponse
	for _, chirp := range chirpsArray {
		chirps = append(chirps, chirpResponse{
			Id:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID.String(),
		})
	}
	err = responseWithJson(w, 200, chirps)
	if err != nil {
		log.Printf("error encoding response :%v", err)
		return
	}
}

func (cfg *ApiConfig) GetChirp(w http.ResponseWriter, r *http.Request) {
	type chirpResponse struct {
		Id        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    string    `json:"user_id"`
	}


	chirpID, err := uuid.Parse(r.PathValue("id")) 
	if err != nil {
		err := responseWithError(w, 400, "Invalid chirp id")
		if err != nil {
			log.Printf("response error :%v", err)
			return
		}
		return
	}
	chirp, err := cfg.DbQueries.GetChirp(r.Context(), chirpID)
	if err != nil && err.Error() == "sql: no rows in result set" {
		err := responseWithError(w, 404, "Chirp not found")
		if err != nil {
			log.Printf("response error :%v", err)
			return
		}
		return
	}
	if err != nil {
		err := responseWithError(w, 500, "Error fetching chirp")
		if err != nil {
			log.Printf("response error :%v", err)
			return
		}
		return
	}

	err = responseWithJson(w, 200, chirpResponse{
		Id:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID.String(),
	})
	if err != nil {
		log.Printf("error encoding response :%v", err)
		return
	}
}



func clearBody(body string) string {
	var restrictedWords []string = []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}
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
