package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/doemoor/web-server/internal/auth"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) DeleteChirp(w http.ResponseWriter, r *http.Request) {

	// read token from header
	headerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Validate the token
	userFromToken, err := auth.ValidateJWT(headerToken, cfg.Secret)
	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			responseWithError(w, http.StatusUnauthorized, "login expired, please login again")
			return
		}
		responseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// read chirp id from request path
	chirpId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		log.Println("DeleteChirp: invalid chirp id: " + err.Error())
		responseWithError(w, 404, "chirp not found")
		return
	}

	// get chirp from db
	chirpFromDB, err := cfg.DbQueries.GetChirp(r.Context(), chirpId)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			responseWithError(w, 404, "chirp not found")
			return
		}
		log.Println("DeleteChirp: error fetching chirp: " + err.Error())
		responseWithError(w, 500, "server error")
		return
	}

	// Check if the user owns the chirp
	if chirpFromDB.UserID != userFromToken {
		responseWithError(w, http.StatusForbidden, "unauthorized")
		return
	}

	// Delete the chirp
	err = cfg.DbQueries.DeleteChirp(r.Context(), chirpId)
	if err != nil {
		log.Println("DeleteChirp: error deleting chirp: " + err.Error())
		responseWithError(w, 500, "server error")
		return
	}

	// Return success
	responseWithJson(w, 204, nil)
}