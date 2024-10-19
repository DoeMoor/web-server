package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/doemoor/web-server/internal/auth"
)

func (cfg *ApiConfig) RevokeUserRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Println("error getting bearer token: " + err.Error())
		responseWithError(w, 401, "Invalid refresh token")
		return
	}

	refreshTokenDB, err := cfg.DbQueries.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			responseWithError(w, 401, "Invalid refresh token")
			log.Println("error getting refresh token: " + err.Error())
			return
		}
		log.Println("error getting refresh token: " + err.Error())
		responseWithError(w, 500, "server error")
		return
	}

	if refreshTokenDB.RevokedAt.Valid {
		responseWithError(w, 409, "token is already revoked")
		return
	}

	err = cfg.DbQueries.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		log.Println("error revoking refresh token: " + err.Error())
		responseWithError(w, 500, "server error")
		return
	}

	w.WriteHeader(204)
}
