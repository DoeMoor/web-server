package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/doemoor/web-server/internal/auth"
	"github.com/doemoor/web-server/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) CreateChirp(w http.ResponseWriter, r *http.Request) {

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Println("CreateChirp: could not get token: " + err.Error())
		responseWithError(w, 401, "there is no token")
		return
	}
	userUUID, err := auth.ValidateJWT(tokenString, cfg.Secret)
	if err != nil {
		log.Println("CreateChirp: invalid token: " + err.Error())
		responseWithError(w, 401, "Invalid token")
		return
	}

	type chirpRequest struct {
		Body string `json:"body"`
	}

	type chirpResponse struct {
		Id        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    string    `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	chirpReq := chirpRequest{}
	err = decoder.Decode(&chirpReq)
	if err != nil {
		log.Println("CreateChirp: decoding json failed: " + err.Error())
		responseWithError(w, 400, "Invalid JSON")
		return
	}

	if len(chirpReq.Body) > 140 {
		responseWithError(w, 400, "Chirp too long")
		return
	}

	cleanBody := clearBody(chirpReq.Body)

	reqParams := database.CreateChirpParams{
		Body:   cleanBody,
		UserID: userUUID,
	}

	chirpDb, err := cfg.DbQueries.CreateChirp(r.Context(), reqParams)
	if err != nil {
		log.Println("CreateChirp: error creating chirp: " + err.Error())
		responseWithError(w, 500, "server error")
		return
	}
	responseWithJson(w, 201, chirpResponse{
		Id:        chirpDb.ID.String(),
		CreatedAt: chirpDb.CreatedAt,
		UpdatedAt: chirpDb.UpdatedAt,
		Body:      chirpDb.Body,
		UserId:    chirpDb.UserID.String(),
	})
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
		log.Println("GetChirps: error fetching chirps: " + err.Error())
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
	responseWithJson(w, 200, chirps)
}

func (cfg *ApiConfig) GetChirp(w http.ResponseWriter, r *http.Request) {
	type chirpResponse struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    uuid.UUID `json:"user_id"`
	}

	chirpID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		log.Println("GetChirp: invalid chirp id: " + err.Error())
		responseWithError(w, 400, "Invalid chirp id")
		return
	}

	// get chirp from db
	chirp, err := cfg.DbQueries.GetChirp(r.Context(), chirpID)
	if err != nil && err.Error() == "sql: no rows in result set" {
		responseWithError(w, 404, "Chirp not found")
		return
	}

	if err != nil {
		log.Println("GetChirp: error fetching chirp: " + err.Error())
		responseWithError(w, 500, "server error")
		return
	}

	responseWithJson(w, 200, chirpResponse{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
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
