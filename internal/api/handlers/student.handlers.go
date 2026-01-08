package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/SecureParadise/go_attendence/internal/api/middleware"
	"github.com/SecureParadise/go_attendence/internal/db"
	"github.com/SecureParadise/go_attendence/internal/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type studentHandler struct {
	store db.Store
}

func NewStudentHandler(store db.Store) *studentHandler {
	return &studentHandler{store}
}

type CreateStudentRequest struct {
	RollNo     string `json:"roll_no" binding:"required"`
	FirstName  string `json:"first_name" binding:"required"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name" binding:"required"`
	Image      string `json:"image"`
	Batch      string `json:"batch" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	BranchCode string `json:"branch_code" binding:"required"`
	SemesterNo int32  `json:"semester_no" binding:"required"`
}

// CreateStudent completes student profile
// @Summary Complete student profile
// @Description Complete student profile with personal and academic details
// @Tags students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateStudentRequest true "Student profile data"
// @Success 201 {object} sqlc.Student
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /student_reg [post]
func (h *studentHandler) CreateStudent(ctx *gin.Context) {
	var req CreateStudentRequest
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

	// 1.1 Check if user is a student
	if user.UserRole != sqlc.UserroleStudent {
		ctx.Error(middleware.NewAPIError(http.StatusBadRequest, "only students can complete student profile", nil))
		return
	}

	// 1.2 Check if profile is already completed
	if user.IsProfileCompleted {
		ctx.Error(middleware.NewAPIError(http.StatusBadRequest, "profile already completed", nil))
		return
	}

	// 2. Fetch branch by code
	branch, err := h.store.GetBranchByCode(ctx, strings.ToUpper(req.BranchCode))
	if err != nil {
		ctx.Error(middleware.NewAPIError(http.StatusNotFound, "branch not found", err))
		return
	}

	// 3. Fetch semester by number and branch
	semArg := sqlc.GetSemesterByNumberAndBranchParams{
		Number:   req.SemesterNo,
		BranchID: branch.ID,
	}
	semester, err := h.store.GetSemesterByNumberAndBranch(ctx, semArg)
	if err != nil {
		ctx.Error(middleware.NewAPIError(http.StatusNotFound, fmt.Sprintf("semester %d not found for branch %s", req.SemesterNo, req.BranchCode), err))
		return
	}

	// 4. Create student and update user profile status in a transaction
	var student sqlc.Student
	err = h.store.WithTx(ctx, func(q *sqlc.Queries) error {
		studentArg := sqlc.CreateStudentParams{
			RollNo:    req.RollNo,
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
			Batch: pgtype.Text{
				String: req.Batch,
				Valid:  true,
			},
			UserID:   user.ID,
			BranchID: branch.ID,
			CurrentSemesterID: pgtype.UUID{
				Bytes: semester.ID,
				Valid: true,
			},
		}

		var err error
		student, err = q.CreateStudent(ctx, studentArg)
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

	ctx.JSON(http.StatusCreated, student)
}

// GetStudentByRollNo returns student details by roll number
// @Summary Get student by roll number
// @Description Fetch student details using their roll number
// @Tags students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param roll_no path string true "Roll Number"
// @Success 200 {object} sqlc.Student
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /student/{roll_no} [get]
func (h *studentHandler) GetStudentByRollNo(ctx *gin.Context) {
	rollNo := ctx.Param("roll_no")
	if rollNo == "" {
		ctx.Error(middleware.NewAPIError(http.StatusBadRequest, "roll_no is required", nil))
		return
	}

	student, err := h.store.GetStudentByRollNo(ctx, rollNo)
	if err != nil {
		ctx.Error(middleware.NewAPIError(http.StatusNotFound, "student not found", err))
		return
	}

	ctx.JSON(http.StatusOK, student)
}
