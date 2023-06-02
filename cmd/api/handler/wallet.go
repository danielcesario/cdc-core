package handler

import (
	"context"
	"net/http"

	"github.com/danielcesario/cdc-core/internal/wallet"
	"github.com/gin-gonic/gin"
)

type WalletService interface {
	Create(ctx context.Context, request wallet.WalletRequest) (*wallet.WalletResponse, error)
	List(ctx context.Context) ([]*wallet.WalletResponse, error)
	AddCollaborator(ctx context.Context, walletCode string, request wallet.WalletCollaboratorRequest) error
	GetByCode(ctx context.Context, code string) (*wallet.WalletResponse, error)
}

type WalletHandler struct {
	service WalletService
}

func NewWalletHandler(service WalletService) *WalletHandler {
	return &WalletHandler{
		service: service,
	}
}

func (h *Handler) CreateWallet(context *gin.Context) {
	var request wallet.WalletRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	response, err := h.wallet.service.Create(context, request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusCreated, response)
}

func (h *Handler) ListWallets(context *gin.Context) {
	response, err := h.wallet.service.List(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, response)
}

func (h *Handler) AddCollaborator(context *gin.Context) {
	walletCode := context.Param("walletCode")

	var request wallet.WalletCollaboratorRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	err := h.wallet.service.AddCollaborator(context, walletCode, request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.Status(http.StatusOK)
}

func (h *Handler) GetWallet(context *gin.Context) {
	walletCode := context.Param("walletCode")

	response, err := h.wallet.service.GetByCode(context, walletCode)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, response)
}
