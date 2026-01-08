package handlers

import (
	"net/http"
	"strings"

	"github.com/SecureParadise/go_attendence/internal/api/middleware"
	"github.com/SecureParadise/go_attendence/internal/db"
	"github.com/SecureParadise/go_attendence/internal/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type teacherHandler struct {
	store db.Store
}

func NewTeacherHandler(store db.Store) *teacherHandler {
	return &teacherHandler{store}
}

type CreateTeacherRequest struct {
	CardNo         string `json:"card_no" binding:"required"`
	FirstName      string `json:"first_name" binding:"required"`
	MiddleName     string `json:"middle_name"`
	LastName       string `json:"last_name" binding:"required"`
	Image          string `json:"image"`
	Email          string `json:"email" binding:"required,email"`
	DepartmentName string `json:"department_name" binding:"required"`
}

// CreateTeacher completes teacher profile
// @Summary Complete teacher profile
// @Description Complete teacher profile with personal and academic details
// @Tags teachers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateTeacherRequest true "Teacher profile data"
// @Success 201 {object} sqlc.Teacher
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /teacher_reg [post]
func (h *teacherHandler) CreateTeacher(ctx *gin.Context) {
	var req CreateTeacherRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	// 1. Fetch user by email
	user, err := h.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		ctx.Error(middleware.NewAPIError(http.StatusNotFound, "user not found", err))
		return
	}

	// 1.1 Check if user is a teacher
	if user.UserRole != sqlc.UserroleTeacher {
		ctx.Error(middleware.NewAPIError(http.StatusBadRequest, "only teachers can complete teacher profile", nil))
		return
	}

	// 1.2 Check if profile is already completed
	if user.IsProfileCompleted {
		ctx.Error(middleware.NewAPIError(http.StatusBadRequest, "profile already completed", nil))
		return
	}

	// 2. Fetch department by name
	dept, err := h.store.GetDepartmentByName(ctx, strings.ToLower(req.DepartmentName))
	if err != nil {
		ctx.Error(middleware.NewAPIError(http.StatusNotFound, "department not found", err))
		return
	}

	// 3. Create teacher and update user profile status in a transaction
	var teacher sqlc.Teacher
	err = h.store.WithTx(ctx, func(q *sqlc.Queries) error {
		teacherArg := sqlc.CreateTeacherParams{
			CardNo:    req.CardNo,
			FirstName: req.FirstName,
			MiddleName: pgtype.Text{
				String: req.MiddleName,
				Valid:  req.MiddleName != "",
			},
			LastName: req.LastName,
			Image: pgtype.Text{
				String: req.Image,
				Valid:  req.Image != "",
			},
			UserID:       user.ID,
			DepartmentID: dept.ID,
		}

		var err error
		teacher, err = q.CreateTeacher(ctx, teacherArg)
		if err != nil {
			return err
		}

		_, err = q.UpdateUserProfileCompleted(ctx, sqlc.UpdateUserProfileCompletedParams{
			ID:                 user.ID,
			IsProfileCompleted: true,
		})
		return err
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, teacher)
}

// GetTeacherByCardNo returns teacher details by card number
// @Summary Get teacher by card number
// @Description Fetch teacher details using their card number
// @Tags teachers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param card_no path string true "Card Number"
// @Success 200 {object} sqlc.Teacher
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /teacher/{card_no} [get]
func (h *teacherHandler) GetTeacherByCardNo(ctx *gin.Context) {
	cardNo := ctx.Param("card_no")
	if cardNo == "" {
		ctx.Error(middleware.NewAPIError(http.StatusBadRequest, "card_no is required", nil))
		return
	}

	teacher, err := h.store.GetTeacherByCardNo(ctx, cardNo)
	if err != nil {
		ctx.Error(middleware.NewAPIError(http.StatusNotFound, "teacher not found", err))
		return
	}

	ctx.JSON(http.StatusOK, teacher)
}
