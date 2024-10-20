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

func (cfg *ApiConfig) Login(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type loginResponse struct {
		Id           string    `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
		IsChirpyRed  bool      `json:"is_chirpy_red"`
	}

	// read request body
	userReq := loginRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userReq)
	if err != nil {
		log.Println("Login: invalid JSON: " + err.Error())
		responseWithError(w, 409, "Invalid JSON")
		return
	}

	// Check if user exists in DB
	userFromDB, err := cfg.DbQueries.GetUserByEmail(r.Context(), userReq.Email)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			responseWithError(w, 401, "Incorrect email or password")
			return
		}
		log.Println("Login: error fetching user: " + err.Error())
		responseWithError(w, 500, "server error")
		return
	}

	// Check if password is correct
	if err := auth.CheckPasswordHash(userReq.Password, userFromDB.HashedPassword); err != nil {
		responseWithError(w, 401, "Incorrect email or password")
		return
	}

	//set token expire in 1 hour
	tokenExpiry := time.Hour * 1

	// Generate token
	tokenString, err := auth.MakeJWT(userFromDB.ID, cfg.Secret, tokenExpiry)
	if err != nil {
		log.Println("Login: error generating token: " + err.Error())
		responseWithError(w, 500, "Error generating token")
		return
	}

	// Generate refresh token
	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		log.Println("Login: error generating refresh token: " + err.Error())
		responseWithError(w, 500, "Error generating refresh token")
		return
	}

	// Save refresh token in DB
	crtParams := database.CreateRefreshTokenParams{
		Token:     refreshTokenString,
		UserID:    userFromDB.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60).UTC(),
	}

	tokenDBString, err := cfg.DbQueries.CreateRefreshToken(r.Context(), crtParams)
	if err != nil {
		log.Println("Login: error saving refresh token: " + err.Error())
		responseWithError(w, 500, "server error")
		return
	}

	// Create response
	userRes := loginResponse{
		Id:           userFromDB.ID.String(),
		CreatedAt:    userFromDB.CreatedAt,
		UpdatedAt:    userFromDB.UpdatedAt,
		Email:        userFromDB.Email,
		Token:        tokenString,
		RefreshToken: tokenDBString,
		IsChirpyRed:  userFromDB.IsChirpyRed,
	}
	responseWithJson(w, 200, userRes)
}
