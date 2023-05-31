package handler

import (
	"context"
	"net/http"

	"github.com/danielcesario/cdc-core/internal/auth"
	"github.com/danielcesario/cdc-core/internal/user"
	"github.com/gin-gonic/gin"
)

type UserService interface {
	Register(ctx context.Context, request user.UserRequest) (*user.UserResponse, error)
	CheckEmailAndPassword(ctx context.Context, email, password string) (*user.User, error)
	ListUsers(ctx context.Context) ([]*user.UserResponse, error)
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *Handler) Register(context *gin.Context) {
	var request user.UserRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if err := request.HashPassword(request.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	response, err := h.user.service.Register(context, request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusCreated, response)
}

func (h *Handler) Login(context *gin.Context) {
	type TokenRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var request TokenRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	user, err := h.user.service.CheckEmailAndPassword(context, request.Email, request.Password)
	if err != nil {
		if err.Error() == "invalid credentials" {
			context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		context.Abort()
		return
	}

	token, err := auth.GenerateJWT(user.Email, user.Code, user.GetRoles())
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) GetUsers(context *gin.Context) {
	response, err := h.user.service.ListUsers(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, response)
}
