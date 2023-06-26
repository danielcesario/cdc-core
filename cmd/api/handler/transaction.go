package handler

import (
	"context"
	"net/http"

	"github.com/danielcesario/cdc-core/internal/transaction"
	"github.com/gin-gonic/gin"
)

type TransactionService interface {
	Create(ctx context.Context, request transaction.TransactionRequest) (*transaction.TransactionResponse, error)
	GetByCode(ctx context.Context, code string) (*transaction.TransactionResponse, error)
	Search(ctx context.Context, request transaction.SearchTransactionRequest) (*transaction.SearchPageResponse, error)
}

type TransactionHandler struct {
	service TransactionService
}

func NewTransactionHandler(service TransactionService) *TransactionHandler {
	return &TransactionHandler{
		service: service,
	}
}

func (h *Handler) CreateTransaction(context *gin.Context) {
	var request transaction.TransactionRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	response, err := h.transaction.service.Create(context, request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusCreated, response)
}

func (h *Handler) GetTransaction(context *gin.Context) {
	transactionCode := context.Param("transactionCode")

	response, err := h.transaction.service.GetByCode(context, transactionCode)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, response)
}

func (h *Handler) SerachTransactions(context *gin.Context) {
	var request transaction.SearchTransactionRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	response, err := h.transaction.service.Search(context, request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, response)
}
