package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"simple-finance/internal/db"
	"simple-finance/internal/handler/middleware"
	"simple-finance/internal/handler/response"
	"simple-finance/internal/models"
	"simple-finance/internal/tokens"
)

type TransactionHandler struct {
	db        *db.FinanceDB
	validator *validator.Validate
	logger    *logrus.Logger
}

func NewTransactionHandler(
	db *db.FinanceDB,
	validator *validator.Validate,
	logger *logrus.Logger,
) *TransactionHandler {
	return &TransactionHandler{
		db:        db,
		validator: validator,
		logger:    logger,
	}
}

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
