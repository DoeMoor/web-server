package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/doemoor/wed-server/internal/auth"
)

func (cfg *ApiConfig) Login(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	userReq := loginRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userReq)
	if err != nil {
		responseWithError(w, 409, "Invalid JSON")
		return
	}

	userHash, err := cfg.DbQueries.GetUserByEmail(r.Context(), userReq.Email)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			responseWithError(w, 401, "Incorrect email or password")
			return
		}
		responseWithError(w, 500, "Error fetching user")
		return
	}

	if err := auth.CheckPasswordHash(userReq.Password, userHash); err != nil {
		responseWithError(w, 401, "Incorrect email or password")
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte("OK"))
}
