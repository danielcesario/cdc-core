package handler

import (
	"context"
	"net/http"

	paymentmethod "github.com/danielcesario/cdc-core/internal/payment-method"
	"github.com/gin-gonic/gin"
)

type PaymentMethodService interface {
	Create(ctx context.Context, request paymentmethod.PaymentMethodRequest) (*paymentmethod.PaymentMethodResponse, error)
	List(ctx context.Context) ([]*paymentmethod.PaymentMethodResponse, error)
}

type PaymentMethodHandler struct {
	service PaymentMethodService
}

func NewPaymentMethodHandler(service PaymentMethodService) *PaymentMethodHandler {
	return &PaymentMethodHandler{
		service: service,
	}
}

func (h *Handler) CreatePaymentMethod(context *gin.Context) {
	var request paymentmethod.PaymentMethodRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	response, err := h.paymentMethod.service.Create(context, request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusCreated, response)
}

func (h *Handler) ListPaymentMethod(context *gin.Context) {
	response, err := h.paymentMethod.service.List(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, response)
}
