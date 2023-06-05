package handler

import (
	"context"
	"net/http"

	category "github.com/danielcesario/cdc-core/internal/category"
	"github.com/gin-gonic/gin"
)

type CategoryService interface {
	Create(ctx context.Context, request category.CategoryRequest) (*category.CategoryResponse, error)
	List(ctx context.Context) ([]*category.CategoryResponse, error)
	Update(ctx context.Context, code string, request category.CategoryRequest) (*category.CategoryResponse, error)
}

type CategoryHandler struct {
	service CategoryService
}

func NewCategoryHandler(service CategoryService) *CategoryHandler {
	return &CategoryHandler{
		service: service,
	}
}

func (h *Handler) CreateCategory(context *gin.Context) {
	var request category.CategoryRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	response, err := h.category.service.Create(context, request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusCreated, response)
}

func (h *Handler) ListCategory(context *gin.Context) {
	response, err := h.category.service.List(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, response)
}

func (h *Handler) UpdateCategory(context *gin.Context) {
	var request category.CategoryRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	categoryCode := context.Param("categoryCode")

	response, err := h.category.service.Update(context, categoryCode, request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusCreated, response)
}
