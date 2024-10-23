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

type UserResponse struct {
	ID           uuid.UUID `json:"id,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
	Email        string    `json:"email,omitempty"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type UserRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (cfg *ApiConfig) HandlerCreateUsers(w http.ResponseWriter, r *http.Request) {

	body := &UserRequest{}
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		utils.RespondWithError(w, 403, "Something went wrong")
		return
	}

	hashedPassword, err := auth.HashPassword(body.Password)

	if err != nil {
		utils.RespondWithError(w, 403, "Something went wrong")
		return
	}

	newUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          body.Email,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		utils.RespondWithError(w, 500, err.Error())
		return
	}

	userResponse := UserResponse{
		ID:        newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email:     newUser.Email,
	}

	utils.RespondWithJSON(w, 201, 500, userResponse)
}

func (cfg *ApiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	req := &UserRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, 403, err.Error())
		return
	}

	user, err := cfg.db.GetUser(r.Context(), req.Email)

	if err != nil {
		utils.RespondWithError(w, 403, err.Error())
		return
	}

	if err := auth.CheckPasswordHash(req.Password, user.HashedPassword); err != nil {
		utils.RespondWithError(w, 401, "Incorrect email or password")
		return
	}

	jwtToken, err := auth.MakeJWT(user.ID, cfg.secretToken, time.Duration(3600)*time.Second)

	if err != nil {
		utils.RespondWithError(w, 500, err.Error())
		return
	}

	refreshToken, err := auth.MakeRefreshToken()

	if err != nil {
		utils.RespondWithError(w, 500, err.Error())
		return
	}

	_, err = cfg.db.SaveRefreshToken(r.Context(), database.SaveRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Duration(3600*24*60) * time.Second),
	})

	if err != nil {
		utils.RespondWithError(w, 500, err.Error())
		return
	}

	userResponse := UserResponse{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        jwtToken,
		RefreshToken: refreshToken,
	}

	utils.RespondWithJSON(w, 200, 500, userResponse)
}

func (cfg *ApiConfig) HandlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)

	if err != nil {
		utils.RespondWithError(w, 403, err.Error())
		return
	}

	refreshTokenData, err := cfg.db.GetRefreshToken(r.Context(), refreshToken)

	if err != nil {
		utils.RespondWithError(w, 401, err.Error())
		return
	}

	utils.ValidateRefreshToken(refreshTokenData.RevokedAt.Valid, refreshTokenData.ExpiresAt, w)

	user, err := cfg.db.GetUserByRefreshToken(r.Context(), refreshToken)

	if err != nil {
		utils.RespondWithError(w, 500, err.Error())
		return
	}

	jwtToken, err := auth.MakeJWT(user.ID, cfg.secretToken, time.Duration(3600)*time.Second)

	if err != nil {
		utils.RespondWithError(w, 500, err.Error())
		return
	}

	userResponse := TokenResponse{
		Token: jwtToken,
	}

	utils.RespondWithJSON(w, 200, 500, userResponse)
}

func (cfg *ApiConfig) HandlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)

	if err != nil {
		utils.RespondWithError(w, 403, err.Error())
		return
	}

	err = cfg.db.RevokeToken(r.Context(), refreshToken)

	if err != nil {
		utils.RespondWithError(w, 500, err.Error())
		return
	}

	w.WriteHeader(204)
}

func (cfg *ApiConfig) HandlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)

	if err != nil {
		utils.RespondWithError(w, 401, err.Error())
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, cfg.secretToken)

	if err != nil {
		utils.RespondWithError(w, 401, err.Error())
		return
	}

	req := &UserRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, 403, err.Error())
		return
	}

	if len(req.Email) == 0 || len(req.Password) == 0 {
		utils.RespondWithError(w, 400, "Email and password are required")
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)

	if err != nil {
		utils.RespondWithError(w, 403, "Something went wrong")
		return
	}

	err = cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		utils.RespondWithError(w, 500, err.Error())
		return
	}

	userResponse := UserResponse{
		ID:        userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     req.Email,
	}

	utils.RespondWithJSON(w, 200, 500, userResponse)
}
