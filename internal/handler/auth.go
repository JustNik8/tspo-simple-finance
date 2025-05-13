package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	_ "net/http"
	"simple-finance/internal/auth"
	"simple-finance/internal/db"
	"simple-finance/internal/errs"
	"simple-finance/internal/handler/response"
	"simple-finance/internal/models"
	"simple-finance/pkg/hash"
)

type AuthHandler struct {
	validate        *validator.Validate
	db              *db.FinanceDB
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	logger          *logrus.Logger
	hasher          hash.PasswordHasher
	authManager     *auth.Manager
}

func NewAuthHandler(
	validate *validator.Validate,
	db *db.FinanceDB,
	logger *logrus.Logger,
	hasher hash.PasswordHasher,
	authManager *auth.Manager,
) *AuthHandler {
	return &AuthHandler{
		validate:        validate,
		db:              db,
		accessTokenTTL:  30 * 24 * time.Hour,
		refreshTokenTTL: 60 * 24 * time.Hour,
		logger:          logger,
		hasher:          hasher,
		authManager:     authManager,
	}
}

// SignIn             godoc
// @Summary      Authenticate user
// @Description  Login with username and password to get access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body  models.SignInInput  true  "User credentials"
// @Success      200    {object}  models.Tokens
// @Failure      400    {object}  string
// @Failure      401    {object}  string
// @Failure      500    {object}  string
// @Router       /auth/sign_in [post]
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

	userID, err := h.authManager.ComparePassword(ctx, input.Username, input.Password)
	if err != nil {
		h.logger.Warn(err)
		if errors.Is(err, errs.ErrInvalidPassword) {
			response.Unauthorized(w)
			return
		}

		response.InternalServerError(w)
		return
	}

	accessToken, refreshToken, err := h.authManager.MakeTokens(userID, h.accessTokenTTL, h.refreshTokenTTL)
	if err != nil {
		h.logger.Warn(err)
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
		h.logger.Warn(err)
		response.InternalServerError(w)
		return
	}

	response.WriteResponse(w, http.StatusOK, ansBytes)
}

// SignUp             godoc
// @Summary      Регистрирует нового пользователя
// @Description  Create a new user account with email, username and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body  models.SignUpInput  true  "User registration data"
// @Success      201    {object}  models.UserInfo
// @Failure      400    {object}  string
// @Failure      500    {object}  string
// @Router       /auth/sign_up [post]
func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var input models.SignUpInput

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

	hashPass, err := h.hasher.Hash(input.Password)
	if err != nil {
		h.logger.Warn(err)
		response.InternalServerError(w)
		return
	}

	ctx := context.Background()
	userInfo, err := h.db.InsertUser(ctx, models.UserInfo{
		ID:       uuid.New().String(),
		Email:    input.Email,
		UserName: input.UserName,
		Password: hashPass,
	})

	if err != nil {
		h.logger.Warn(err)
		response.InternalServerError(w)
		return
	}

	ansBytes, err := json.Marshal(userInfo)
	if err != nil {
		h.logger.Warn(err)
		response.InternalServerError(w)
		return
	}

	response.WriteResponse(w, http.StatusOK, ansBytes)
}

func (h *AuthHandler) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	var input models.RefreshInput

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

	accessToken, refreshToken, err := h.authManager.RefreshTokens(
		input.RefreshToken,
		h.accessTokenTTL,
		h.refreshTokenTTL,
	)
	if err != nil {
		h.logger.Warn(err)
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
		h.logger.Warn(err)
		response.InternalServerError(w)
		return
	}

	response.WriteResponse(w, http.StatusOK, ansBytes)
}
