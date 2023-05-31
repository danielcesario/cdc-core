package handler

import (
	"context"
	"net/http"

	"github.com/danielcesario/cdc-core/internal/plan"
	"github.com/gin-gonic/gin"
)

type PlanService interface {
	Create(ctx context.Context, request plan.PlanRequest) (*plan.PlanResponse, error)
	ListAll(ctx context.Context) ([]*plan.PlanResponse, error)
	GetByCode(ctx context.Context, code string) (*plan.PlanResponse, error)
	Update(ctx context.Context, planCode string, request plan.PlanRequest) (*plan.PlanResponse, error)
}

type PlanHandler struct {
	service PlanService
}

func NewPlanHandler(service PlanService) *PlanHandler {
	return &PlanHandler{
		service: service,
	}
}

func (h *Handler) CreatePlan(context *gin.Context) {
	var request plan.PlanRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	response, err := h.plan.service.Create(context, request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusCreated, response)
}

func (h *Handler) ListPlans(context *gin.Context) {
	response, err := h.plan.service.ListAll(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, response)
}

func (h *Handler) GetPlan(context *gin.Context) {
	planCode := context.Param("planCode")
	response, err := h.plan.service.GetByCode(context, planCode)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, response)
}

func (h *Handler) UpdatePlan(context *gin.Context) {
	planCode := context.Param("planCode")
	var request plan.PlanRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	response, err := h.plan.service.Update(context, planCode, request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, response)
}
