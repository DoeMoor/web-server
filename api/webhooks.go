package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/doemoor/web-server/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) PolkaWebhook(w http.ResponseWriter, r *http.Request) {

	type webhookRequest struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	webhookReq := webhookRequest{}
	err := decoder.Decode(&webhookReq)
	if err != nil {
		log.Println("Webhooks: could not decode request: " + err.Error())
		responseWithError(w, 400, "could not decode request")
		return
	}

	if webhookReq.Event != "user.upgraded" {
		log.Println("Webhooks: unknown event: " + webhookReq.Event)
		responseWithError(w, 204, "unknown event")
		return
	}

	if webhookReq.Data.UserID == uuid.Nil {
		log.Println("Webhooks: missing user id")
		responseWithError(w, 204, "user id is missing or invalid")
		return
	}

	upbUserArg := database.UpdateUserMembershipChirpyRedParams{
		ID:          webhookReq.Data.UserID,
		IsChirpyRed: true,
	}
	user_id, err := cfg.DbQueries.UpdateUserMembershipChirpyRed(r.Context(), upbUserArg)
	if user_id == uuid.Nil {
		log.Println("Webhooks: could not update user " + user_id.String() + " membership: " + err.Error())
		responseWithError(w, 404, "user not found")
	}
	if err != nil {
		log.Println("Webhooks: could not update user membership: " + err.Error())
		responseWithError(w, 500, "could not update user membership")
		return
	}

	w.WriteHeader(204)
}
