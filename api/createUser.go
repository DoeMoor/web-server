package api

import (
	"encoding/json"
	"net/http"
	"time"

	// "github.com/doemoor/wed-server/internal/database"
)

func (cfg *ApiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	type userRequest struct {
		Email string `json:"email"`
	}

	type userResponse struct {
		Id        string `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string `json:"email"`
	}

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	userReq := userRequest{}
	err := decoder.Decode(&userReq)
	if err != nil {
		responseWithError(w, 500, "Invalid JSON")
		return
	}

	userDb, err := cfg.DbQueries.CreateUser(r.Context(), userReq.Email)
	if err != nil {
		responseWithError(w, 500, "Error creating user")
		return
	}

	userResp := userResponse{
		Id:        userDb.ID.String(),
		CreatedAt: userDb.CreatedAt,
		UpdatedAt: userDb.UpdatedAt,
		Email:     userDb.Email,
	}

	responseWithJson(w, 201, userResp)
}
