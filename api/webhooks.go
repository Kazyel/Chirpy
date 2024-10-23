package api

import (
	"encoding/json"
	"net/http"

	"github.com/Kazyel/chirpy-bootdev/internal/auth"
	"github.com/Kazyel/chirpy-bootdev/utils"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) HandlerUpdateChirpy(w http.ResponseWriter, r *http.Request) {
	type updateRequest struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		}
	}

	apiKey, err := auth.GetAPIKey(r.Header)

	if err != nil {
		utils.RespondWithError(w, 401, err.Error())
		return
	}

	if apiKey != cfg.polkaKey {
		utils.RespondWithError(w, 401, "Invalid API key")
		return
	}

	req := &updateRequest{}
	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		utils.RespondWithError(w, 403, err.Error())
		return
	}

	if req.Event != "user.upgraded" {
		w.WriteHeader(204)
	}

	err = cfg.db.UpgradeChirpy(r.Context(), req.Data.UserID)

	if err != nil {
		utils.RespondWithError(w, 404, err.Error())
	}

	utils.RespondWithError(w, 204, "")
}
