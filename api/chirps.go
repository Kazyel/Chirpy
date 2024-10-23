package api

import (
	"encoding/json"
	"net/http"
	"time"

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
		Body string `json:"body"`
	}

	req := &chirpRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		utils.RespondWithError(w, 403, "Something went wrong")
		return
	}

	if len(req.Body) > 140 {
		utils.RespondWithError(w, 400, "Chirpy is too long")
	}

	filteredBody := utils.ProfaneFilter(req.Body)
	userID, err := uuid.Parse(r.Header.Get("User-ID"))

	if err != nil || userID == uuid.Nil {
		utils.RespondWithError(w, 401, "Invalid ID")
		return
	}

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

	utils.RespondWithJSON(w, 201, 500, chirpResponse)
}

func (cfg *ApiConfig) HandlerGetChirps(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author_id")

	if authorID != "" {
		uuid, err := uuid.Parse(authorID)

		if err != nil {
			utils.RespondWithError(w, 400, "Invalid ID")
			return
		}

		data, err := cfg.db.GetChirpsByAuthor(r.Context(), uuid)

		if err != nil {
			utils.RespondWithError(w, 404, err.Error())
		}

		chirpsResponse := make([]ChirpResponse, len(data))

		for index, chirp := range data {
			chirpsResponse[index] = ChirpResponse{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserID:    chirp.UserID,
			}
		}

		utils.RespondWithJSON(w, 200, 500, chirpsResponse)
	}

	data, err := cfg.db.GetChirps(r.Context())

	if err != nil {
		utils.RespondWithError(w, 500, err.Error())
		return
	}

	chirpsResponse := make([]ChirpResponse, len(data))

	for index, chirp := range data {
		chirpsResponse[index] = ChirpResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	utils.RespondWithJSON(w, 200, 500, chirpsResponse)
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

	chirpResponse := ChirpResponse{
		ID:        data.ID,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		Body:      data.Body,
		UserID:    data.UserID,
	}

	utils.RespondWithJSON(w, 200, 500, chirpResponse)
}

func (cfg *ApiConfig) HandlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.Header.Get("User-ID"))

	if err != nil || userID == uuid.Nil {
		utils.RespondWithError(w, 401, "Invalid ID")
		return
	}

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))

	if err != nil {
		utils.RespondWithError(w, 500, err.Error())
		return
	}

	_, err = cfg.db.GetChirpsByID(r.Context(), chirpID)

	if err != nil {
		utils.RespondWithError(w, 404, err.Error())
		return
	}

	deletedChirp, err := cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     chirpID,
		UserID: userID,
	})

	if err != nil {
		utils.RespondWithError(w, 403, err.Error())
		return
	}

	if (deletedChirp == database.Chirp{}) {
		utils.RespondWithError(w, 404, "Chirp not found")
		return
	}

	w.WriteHeader(204)
}
