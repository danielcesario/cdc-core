package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	user          *UserHandler
	plan          *PlanHandler
	wallet        *WalletHandler
	paymentMethod *PaymentMethodHandler
}

func NewHandler(user *UserHandler, plan *PlanHandler, wallet *WalletHandler, paymentMethod *PaymentMethodHandler) *Handler {
	return &Handler{
		user:          user,
		plan:          plan,
		wallet:        wallet,
		paymentMethod: paymentMethod,
	}
}

func (h *Handler) FakeResponse(context *gin.Context) {
	context.JSON(http.StatusCreated, map[string]interface{}{"result": "created"})
}
