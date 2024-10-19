package api

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/doemoor/web-server/internal/auth"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) RefreshUserToken(w http.ResponseWriter, r *http.Request) {

	// refresh token is provided in header
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Println("error getting bearer token", err)
		responseWithError(w, 401, "Invalid refresh token")
		return
	}

	// get refresh token from db
	refreshTokenDB, err := cfg.DbQueries.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			responseWithError(w, 401, "Invalid refresh token")
			return
		}

		log.Println("error getting refresh token", err)
		responseWithError(w, 500, "server error")
		return
	}

	// check if refresh token is NOT revoked!!
	if refreshTokenDB.RevokedAt.Valid {
		responseWithError(w, 401, "login expired, please login again")
		return
	}

	// check if refresh token is expired
	if refreshTokenDB.ExpiresAt.Before(time.Now().UTC()) {
		err := cfg.DbQueries.RevokeRefreshToken(r.Context(), refreshTokenDB.Token)
		if err != nil {
			log.Println("error revoking refresh token" + err.Error())
			responseWithError(w, 500, "server error")
			return
		}
		responseWithError(w, 401, "login expired, please login again")
		return
	}

	// generate new token
	newUserToken, err := auth.MakeJWT(uuid.UUID(refreshTokenDB.UserID), cfg.Secret, time.Hour*1)
	if err != nil {
		log.Println("error generating token newUserToken" + err.Error())
		responseWithError(w, 500, "server error")
		return
	}

	type refreshTokenResponse struct {
		Token string `json:"token"`
	}

	// response
	responseWithJson(w, 200, refreshTokenResponse{Token: newUserToken})
}
