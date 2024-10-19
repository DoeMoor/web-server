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


func (cfg *ApiConfig) UpdateUser(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	type userRequest struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	// read request body
	decoder := json.NewDecoder(r.Body)
	userReq := userRequest{}
	err = decoder.Decode(&userReq)
	if err != nil {
		log.Println("UpdateUser: invalid JSON: " + err.Error())
		responseWithError(w, 400, "Invalid JSON")
		return
	}

	// Validate the request
	if userReq.Password == "" || userReq.Email == "" {
		responseWithError(w, 400, "missing password or email")
		return
	}

	// Validate the token 
	userFromToken, err := auth.ValidateJWT(token, cfg.Secret)
	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			responseWithError(w, http.StatusUnauthorized, "login expired, please login again")
			return
		}
		log.Println("UpdateUser: invalid token: " + err.Error())
		responseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get the user from the database
  if ok, _ := cfg.DbQueries.IsUserIdExists(r.Context(), userFromToken); !ok {
		responseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	newHashedPassword, err := auth.HashPassword(userReq.Password)
	if err != nil {
		log.Println("UpdateUser: error hashing password: " + err.Error())
		responseWithError(w, 500, "server error")
		return
	}

	// Update the user in the database
	updateUserArgs := database.UpdateUserParams{
		HashedPassword: newHashedPassword,
		Email:    userReq.Email,
		ID:       userFromToken,
	}

	updatedUser, err := cfg.DbQueries.UpdateUser(r.Context(), updateUserArgs) 
	if err != nil {
		log.Println("UpdateUser: error updating user: " + err.Error())
		responseWithError(w, 500, "server error")
		return
	}

	// Return the updated user
	type userResponse struct {
		Id        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	userResp := userResponse{
		Id:        updatedUser.ID.String(),
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email:     updatedUser.Email,
	}

	responseWithJson(w, 200, userResp)
}
