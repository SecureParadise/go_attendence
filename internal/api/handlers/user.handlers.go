package handlers

import (
	"net/http"

	"github.com/SecureParadise/go_attendence/internal/api/middleware"
	"github.com/SecureParadise/go_attendence/internal/auth"
	"github.com/SecureParadise/go_attendence/internal/config"
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
	store      db.Store
	tokenMaker auth.Maker
	config     config.Config
}

func NewUserHandler(store db.Store, tokenMaker auth.Maker, config config.Config) *userHandler {
	return &userHandler{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}
}

// CreateUser handles user registration
// @Summary Create a new user
// @Description Register a new user with email, password, and optional role
// @Tags users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "User registration data"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
func (h *userHandler) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.Error(middleware.NewAPIError(http.StatusInternalServerError, "failed to hash password", err))
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
		ctx.Error(err)
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
	AccessToken        string `json:"access_token"`
	Email              string `json:"email"`
	Role               string `json:"role"`
	IsProfileCompleted bool   `json:"is_profile_completed"`
}

// Login handles user authentication
// @Summary User login
// @Description Authenticate user and return access token
// @Tags users
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 401 {object} map[string]string
// @Router /login [post]
func (h *userHandler) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	user, err := h.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		ctx.Error(middleware.NewAPIError(http.StatusUnauthorized, "invalid credentials", err))
		return
	}

	err = util.CheckPassword(user.PasswordHash, req.Password)
	if err != nil {
		ctx.Error(middleware.NewAPIError(http.StatusUnauthorized, "invalid credentials", err))
		return
	}

	accessToken, _, err := h.tokenMaker.CreateToken(
		user.Email,
		string(user.UserRole),
		h.config.AccessTokenDuration,
		auth.AccessToken,
	)
	if err != nil {
		ctx.Error(middleware.NewAPIError(http.StatusInternalServerError, "failed to create access token", err))
		return
	}

	rsp := LoginResponse{
		Email:              user.Email,
		Role:               string(user.UserRole),
		IsProfileCompleted: user.IsProfileCompleted,
		AccessToken:        accessToken,
	}

	ctx.JSON(http.StatusOK, rsp)
}

// GetUserMe returns the current authenticated user's profile
// @Summary Get current user profile
// @Description Get profile of the user identified by the access token
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} sqlc.User
// @Failure 401 {object} map[string]string
// @Router /user/me [get]
func (h *userHandler) GetUserMe(ctx *gin.Context) {
	payload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*auth.Payload)

	user, err := h.store.GetUserByEmail(ctx, payload.Username)
	if err != nil {
		ctx.Error(middleware.NewAPIError(http.StatusNotFound, "user not found", err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}
