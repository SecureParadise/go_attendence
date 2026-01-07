package handlers

import (
	"net/http"

	"github.com/SecureParadise/go_attendence/internal/db"
	"github.com/SecureParadise/go_attendence/internal/db/sqlc"
	"github.com/SecureParadise/go_attendence/internal/util"
	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Email    string         `json:"email" binding:"required,email"`
	Password string         `json:"password" binding:"required,min=6"`
	UserRole *sqlc.Userrole `json:"user_role"`
}

type userHandler struct {
	store db.Store
}

func NewUserHandler(store db.Store) *userHandler {
	return &userHandler{store: store}
}

func (h *userHandler) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Use the provided role or default to student
	userRole := sqlc.UserroleStudent
	if req.UserRole != nil {
		userRole = *req.UserRole
	}

	arg := sqlc.CreateUserParams{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		UserRole:     userRole,
	}

	user, err := h.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rsp := LoginResponse{
		Email:              user.Email,
		IsProfileCompleted: user.IsProfileCompleted,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Email              string `json:"email"`
	IsProfileCompleted bool   `json:"is_profile_completed"`
}

func (h *userHandler) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	err = util.CheckPassword(user.PasswordHash, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	rsp := LoginResponse{
		Email:              user.Email,
		IsProfileCompleted: user.IsProfileCompleted,
	}

	ctx.JSON(http.StatusOK, rsp)
}
