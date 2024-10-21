package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Kazyel/chirpy-bootdev/internal/auth"
	"github.com/Kazyel/chirpy-bootdev/internal/database"
	"github.com/Kazyel/chirpy-bootdev/utils"
	"github.com/google/uuid"
)

type ChirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *ApiConfig) HandlerCreateChirps(w http.ResponseWriter, r *http.Request) {
	type chirpRequest struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	bearerToken, err := auth.GetBearerToken(r.Header)

	if err != nil {
		utils.RespondWithError(w, 403, err.Error())
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, cfg.secretToken)

	if err != nil {
		utils.RespondWithError(w, 401, err.Error())
		return
	}

	req := &chirpRequest{}
	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		utils.RespondWithError(w, 403, "Something went wrong")
		return
	}

	if len(req.Body) > 140 {
		utils.RespondWithError(w, 400, "Chirpy is too long")
	}

	filteredBody := utils.ProfaneFilter(req.Body)

	newChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   filteredBody,
		UserID: userID,
	})

	if err != nil {
		utils.RespondWithError(w, 500, err.Error())
		return
	}

	chirpResponse := ChirpResponse{
		ID:        newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body:      newChirp.Body,
		UserID:    newChirp.UserID,
	}

	marshalResponse, err := json.Marshal(chirpResponse)
	if err != nil {
		utils.RespondWithError(w, 500, err.Error())
		return
	}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(201)
	w.Write(marshalResponse)
}

func (cfg *ApiConfig) HandlerGetChirps(w http.ResponseWriter, r *http.Request) {
	data, err := cfg.db.GetChirps(r.Context())

	if err != nil {
		utils.RespondWithError(w, 500, err.Error())
		return
	}

	chirps := make([]ChirpResponse, len(data))
	for index, chirp := range data {
		chirps[index] = ChirpResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	marshalResponse, err := json.Marshal(chirps)
	if err != nil {
		utils.RespondWithError(w, 500, err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	w.Write(marshalResponse)
}

func (cfg *ApiConfig) HandlerGetChirpByUserID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("userID"))

	if err != nil || id == uuid.Nil {
		utils.RespondWithError(w, 400, "Invalid ID")
		return
	}

	data, err := cfg.db.GetChirpsByUserID(r.Context(), id)

	if err != nil {
		utils.RespondWithError(w, 500, err.Error())
		return
	}

	chirps := make([]ChirpResponse, len(data))
	for index, chirp := range data {
		chirps[index] = ChirpResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	marshalResponse, err := json.Marshal(chirps)
	if err != nil {
		utils.RespondWithError(w, 500, err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	w.Write(marshalResponse)
}

func (cfg *ApiConfig) HandlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("chirpID"))

	if err != nil || id == uuid.Nil {
		utils.RespondWithError(w, 400, "Invalid ID")
		return
	}

	data, err := cfg.db.GetChirpsByID(r.Context(), id)

	if err != nil {
		utils.RespondWithError(w, 404, err.Error())
		return
	}

	chirp := ChirpResponse{
		ID:        data.ID,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		Body:      data.Body,
		UserID:    data.UserID,
	}

	marshalResponse, err := json.Marshal(chirp)
	if err != nil {
		utils.RespondWithError(w, 500, err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	w.Write(marshalResponse)
}
