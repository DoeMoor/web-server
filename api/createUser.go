package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/doemoor/wed-server/internal/auth"
	"github.com/doemoor/wed-server/internal/database"
	// "github.com/doemoor/wed-server/internal/database"
)

func (cfg *ApiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	type userRequest struct {
		Password string `json:"password"`
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
		responseWithError(w, 400, "Invalid JSON")
		return
	}

	HashedPassword, err := auth.HashPassword(userReq.Password)
	if err != nil {
		responseWithError(w, 500, "Error hashing password")
		return
	}

	userParams := database.CreateUserParams{
		HashedPassword: HashedPassword,
		Email: userReq.Email,
	}

	userDb, err := cfg.DbQueries.CreateUser(r.Context(), userParams)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			responseWithError(w, 409, "User already exist: " + userReq.Email )
			return
		}
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
