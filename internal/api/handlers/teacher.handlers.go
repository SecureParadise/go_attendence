package handlers

import (
	"net/http"
	"strings"

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

func (h *teacherHandler) CreateTeacher(ctx *gin.Context) {
	var req CreateTeacherRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Fetch user by email
	user, err := h.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// 1.1 Check if user is a teacher
	if user.UserRole != sqlc.UserroleTeacher {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "only teachers can complete teacher profile"})
		return
	}

	// 1.2 Check if profile is already completed
	if user.IsProfileCompleted {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "profile already completed"})
		return
	}

	// 2. Fetch department by name
	dept, err := h.store.GetDepartmentByName(ctx, strings.ToLower(req.DepartmentName))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "department not found"})
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, teacher)
}

func (h *teacherHandler) GetTeacherByCardNo(ctx *gin.Context) {
	cardNo := ctx.Param("card_no")
	if cardNo == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "card_no is required"})
		return
	}

	teacher, err := h.store.GetTeacherByCardNo(ctx, cardNo)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "teacher not found"})
		return
	}

	ctx.JSON(http.StatusOK, teacher)
}
