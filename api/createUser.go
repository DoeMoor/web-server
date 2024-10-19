package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/doemoor/web-server/internal/auth"
	"github.com/doemoor/web-server/internal/database"
)

func (cfg *ApiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	type userRequest struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type userResponse struct {
		Id        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	// read request body
	decoder := json.NewDecoder(r.Body)
	userReq := userRequest{}
	err := decoder.Decode(&userReq)
	if err != nil {
		log.Println("CreateUser: invalid JSON: " + err.Error())
		responseWithError(w, 400, "Invalid JSON")
		return
	}

	// hash password
	HashedPassword, err := auth.HashPassword(userReq.Password)
	if err != nil {
		log.Println("CreateUser: error hashing password: " + err.Error())
		responseWithError(w, 500, "server error")
		return
	}

	// create user in db
	userParams := database.CreateUserParams{
		HashedPassword: HashedPassword,
		Email:          userReq.Email,
	}

	userDb, err := cfg.DbQueries.CreateUser(r.Context(), userParams)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			responseWithError(w, 409, "User already exist: "+userReq.Email)
			return
		}

		log.Println("CreateUser: error creating user: " + err.Error())
		responseWithError(w, 500, "server error")
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
