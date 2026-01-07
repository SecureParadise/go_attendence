package handlers

import (
	"fmt"
	"net/http"
	"strings"

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

func (h *studentHandler) CreateStudent(ctx *gin.Context) {
	var req CreateStudentRequest
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

	// 1.1 Check if user is a student
	if user.UserRole != sqlc.UserroleStudent {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "only students can complete student profile"})
		return
	}

	// 1.2 Check if profile is already completed
	if user.IsProfileCompleted {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "profile already completed"})
		return
	}

	// 2. Fetch branch by code
	branch, err := h.store.GetBranchByCode(ctx, strings.ToUpper(req.BranchCode))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "branch not found"})
		return
	}

	// 3. Fetch semester by number and branch
	semArg := sqlc.GetSemesterByNumberAndBranchParams{
		Number:   req.SemesterNo,
		BranchID: branch.ID,
	}
	semester, err := h.store.GetSemesterByNumberAndBranch(ctx, semArg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("semester %d not found for branch %s", req.SemesterNo, req.BranchCode)})
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, student)
}

func (h *studentHandler) GetStudentByRollNo(ctx *gin.Context) {
	rollNo := ctx.Param("roll_no")
	if rollNo == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "roll_no is required"})
		return
	}

	student, err := h.store.GetStudentByRollNo(ctx, rollNo)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
		return
	}

	ctx.JSON(http.StatusOK, student)
}
