package handler

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"simple-finance/internal/db"
	"simple-finance/internal/handler/middleware"
	"simple-finance/internal/handler/response"
	"simple-finance/internal/models"
	"simple-finance/internal/tokens"
	"time"
)

type TransactionHandler struct {
	db          *db.FinanceDB
	validator   *validator.Validate
	logger      *logrus.Logger
	redisClient *redis.Client
}

func NewTransactionHandler(
	db *db.FinanceDB,
	validator *validator.Validate,
	logger *logrus.Logger,
	redisClient *redis.Client,
) *TransactionHandler {
	return &TransactionHandler{
		db:          db,
		validator:   validator,
		logger:      logger,
		redisClient: redisClient,
	}
}

// InsertTransaction             godoc
// @Summary      Create a new transaction
// @Description  Add a new financial transaction for the authenticated user
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        input  body  models.Transaction  true  "Transaction data"
// @Success      200    {object}  response.IDResponse
// @Failure      400    {object}  string
// @Failure      401    {object}  string
// @Failure      500    {object}  string
// @Router       /api/transaction [post]
// @Security Bearer
func (h *TransactionHandler) InsertTransaction(w http.ResponseWriter, r *http.Request) {
	tokenInfo, ok := r.Context().Value(middleware.TokenInfoKey).(tokens.TokenInfo)
	if !ok {
		response.InternalServerError(w)
		return
	}

	var transaction models.Transaction
	err := json.NewDecoder(r.Body).Decode(&transaction)

	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	err = h.validator.Struct(transaction)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	transaction.UserID = tokenInfo.UserID
	id := uuid.New().String()
	transaction.ID = id
	ctx := context.Background()
	transactionID, err := h.db.InsertTransaction(ctx, transaction)

	if err != nil {
		h.logger.Warn(err)
		response.InternalServerError(w)
		return
	}

	response.IdResponse(w, transactionID)
}

// GetTransactions             godoc
// @Summary      Get user transactions
// @Description  Retrieve all transactions for the authenticated user
// @Tags         transactions
// @Produce      json
// @Success      200  {array}  models.Transaction
// @Failure      401  {object}  string
// @Failure      500  {object}  string
// @Router       /api/transaction [get]
// @Security     Bearer
func (h *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	tokenInfo, ok := r.Context().Value(middleware.TokenInfoKey).(tokens.TokenInfo)
	if !ok {
		h.logger.Info("Not found tokenInfo")
		response.InternalServerError(w)
		return
	}

	ctx := context.Background()
	transactions, err := h.db.GetTransactions(ctx, tokenInfo.UserID)
	if err != nil {
		h.logger.Warn(err)
		response.InternalServerError(w)
		return
	}

	resp, err := json.Marshal(transactions)
	if err != nil {
		h.logger.Warn(err)
		response.InternalServerError(w)
		return
	}
	response.WriteResponse(w, http.StatusOK, resp)
}

// GetTransactionByID             godoc
// @Summary      Get single transaction
// @Description  Get a specific transaction by its ID
// @Tags         transactions
// @Produce      json
// @Param        transaction_uuid  path  string  true  "Transaction UUID"
// @Success      200  {object}  models.Transaction
// @Failure      400  {object}  string
// @Failure      401  {object}  string
// @Failure      500  {object}  string
// @Router       /api/transaction/{transaction_uuid} [get]
// @Security Bearer
func (h *TransactionHandler) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	tokenInfo, ok := r.Context().Value(middleware.TokenInfoKey).(tokens.TokenInfo)
	if !ok {
		response.InternalServerError(w)
		return
	}

	transactionID := chi.URLParam(r, "transaction_uuid")
	if transactionID == "" {
		response.BadRequest(w, "transaction_uuid is empty")
		return
	}

	transaction, err := h.db.GetTransactionByID(context.Background(), tokenInfo.UserID, transactionID)
	if err != nil {
		h.logger.Warn(err)
		response.InternalServerError(w)
		return
	}

	resp, err := json.Marshal(transaction)
	if err != nil {
		h.logger.Warn(err)
		response.InternalServerError(w)
		return
	}

	response.WriteResponse(w, http.StatusOK, resp)
}

// DeleteTransactionByID             godoc
// @Summary      Delete transaction
// @Description  Delete a specific transaction by its ID
// @Tags         transactions
// @Produce      json
// @Param        transaction_uuid  path  string  true  "Transaction UUID"
// @Success      200
// @Failure      400  {object}  string
// @Failure      401  {object}  string
// @Failure      404  {object}  string
// @Failure      500  {object}  string
// @Router       /api/transaction/{transaction_uuid} [delete]
// @Security     Bearer
func (h *TransactionHandler) DeleteTransactionByID(w http.ResponseWriter, r *http.Request) {
	tokenInfo, ok := r.Context().Value(middleware.TokenInfoKey).(tokens.TokenInfo)
	if !ok {
		response.InternalServerError(w)
		return
	}

	transactionID := chi.URLParam(r, "transaction_uuid")
	if transactionID == "" {
		response.BadRequest(w, "transaction_uuid is empty")
		return
	}

	err := h.db.DeleteTransactionByID(context.Background(), tokenInfo.UserID, transactionID)
	if err != nil {
		h.logger.Warn(err)
		response.InternalServerError(w)
		return
	}

	response.IdResponse(w, transactionID)
}

// GetProfileHandler             godoc
// @Summary      Get profile
// @Description  Get profile by its ID
// @Tags         transactions
// @Produce      json
// @Param        id  path  string  true  "id"
// @Success      200
// @Failure      400  {object}  string
// @Failure      401  {object}  string
// @Failure      404  {object}  string
// @Failure      500  {object}  string
// @Router       /api/profile/{id} [get]
// @Security     Bearer
func (r *TransactionHandler) GetProfileHandler(w http.ResponseWriter, req *http.Request) {
	userID := chi.URLParam(req, "id")
	ctx := req.Context()

	//// Попытка получить данные из кеша
	cachedData, err := r.redisClient.Get(ctx, "profile:"+userID).Bytes()
	if err == nil {
		r.logger.Infof("Got profile from cache: %s", string(cachedData))

		w.Header().Set("Content-Type", "application/json")
		//w.Write([]byte("from redis"))
		json.NewEncoder(w).Encode(cachedData)

		return
	}

	userInfo, err := r.db.GetUserById(ctx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responseData, _ := json.Marshal(userInfo)

	// Сохраняем в кеш на 10 минут
	r.redisClient.Set(ctx, "profile:"+userID, responseData, 10*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseData)
}
