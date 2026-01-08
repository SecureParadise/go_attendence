package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/SecureParadise/go_attendence/internal/api/middleware"
	"github.com/SecureParadise/go_attendence/internal/db"
	"github.com/SecureParadise/go_attendence/internal/db/sqlc"
	"github.com/gin-gonic/gin"
)

// struct to deal with DB
type semesterHandler struct {
	store db.Store
}

func NewSemesterHandler(store db.Store) *semesterHandler {
	return &semesterHandler{store: store}
}

type CreateSemesterRequest struct {
	Number uint8  `json:"semester" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Code   string `json:"code" binding:"required"`
}

func (h *semesterHandler) CreateSemester(ctx *gin.Context) {
	var req CreateSemesterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	sem, err := h.createSemesterInternal(ctx, h.store, req)
	if err != nil {
		if err.Error() == "Branch not found" {
			ctx.Error(middleware.NewAPIError(http.StatusNotFound, "branch not found", err))
		} else {
			ctx.Error(err)
		}
		return
	}
	ctx.JSON(http.StatusCreated, sem)

}

func (h *semesterHandler) createSemesterInternal(ctx *gin.Context, q sqlc.Querier, req CreateSemesterRequest) (sqlc.Semester, error) {
	// convert code name in uppercase
	branch, err := q.GetBranchByCode(ctx, strings.ToUpper(req.Code))

	if err != nil {
		return sqlc.Semester{}, fmt.Errorf("Branch not found")

	}
	arg := sqlc.CreateSemesterParams{
		Number:   int32(req.Number),
		Name:     strings.ToUpper(req.Name),
		BranchID: branch.ID,
	}
	return q.CreateSemester(ctx, arg)

}
