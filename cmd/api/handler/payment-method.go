package handler

import (
	"context"
	"net/http"

	paymentmethod "github.com/danielcesario/cdc-core/internal/payment-method"
	"github.com/gin-gonic/gin"
)

type PaymentMethodService interface {
	CreatePaymentMethod(ctx context.Context, request paymentmethod.PaymentMethodRequest) (*paymentmethod.PaymentMethodResponse, error)
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

	response, err := h.paymentMethod.service.CreatePaymentMethod(context, request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusCreated, response)
}
