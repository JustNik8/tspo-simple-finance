package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"simple-finance/internal/db"
	"simple-finance/internal/handler/response"
	"simple-finance/internal/models"
	"simple-finance/internal/tokens"
)

type AuthHandler struct {
	validate        *validator.Validate
	tokenManager    *tokens.TokenManager
	db              *db.FinanceDB
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAuthHandler(
	validate *validator.Validate,
	tokenManager *tokens.TokenManager,
	db *db.FinanceDB,
) *AuthHandler {
	return &AuthHandler{
		validate:        validate,
		tokenManager:    tokenManager,
		db:              db,
		accessTokenTTL:  30 * 24 * time.Hour,
		refreshTokenTTL: 30 * 24 * time.Hour,
	}
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var input models.SignInInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	err = h.validate.Struct(input)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	ctx := context.Background()

	userID, err := h.db.GetUserID(ctx, input.Username)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	accessToken, err := h.tokenManager.NewJWT(tokens.TokenInfo{UserID: userID}, h.accessTokenTTL)
	if err != nil {
		log.Println(err)
		response.InternalServerError(w)
		return
	}

	refreshToken, err := h.tokenManager.NewJWT(tokens.TokenInfo{UserID: userID}, h.refreshTokenTTL)
	if err != nil {
		log.Println(err)
		response.InternalServerError(w)
		return
	}

	ansBytes, err := json.Marshal(
		models.Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	)
	if err != nil {
		response.InternalServerError(w)
		return
	}

	response.WriteResponse(w, http.StatusOK, ansBytes)
}
