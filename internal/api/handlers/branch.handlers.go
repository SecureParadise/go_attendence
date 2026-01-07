package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/SecureParadise/go_attendence/internal/db"
	"github.com/SecureParadise/go_attendence/internal/db/sqlc"
	"github.com/gin-gonic/gin"
)

type CreateBranchRequest struct {
	Name           string `json:"name" binding:"required"`
	Code           string `json:"code" binding:"required"`
	DepartmentName string `json:"department_name" binding:"required"`
}

type branchHandler struct {
	store db.Store
}

func NewBranchHandler(store db.Store) *branchHandler {
	return &branchHandler{store: store}
}

func (h *branchHandler) CreateBranch(ctx *gin.Context) {
	var req CreateBranchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	branch, err := h.createBranchInternal(ctx, h.store, req)
	if err != nil {
		if err.Error() == "department not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, branch)
}

func (h *branchHandler) BulkCreateBranches(ctx *gin.Context) {
	var req []CreateBranchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var branches []sqlc.Branch
	err := h.store.WithTx(ctx, func(q *sqlc.Queries) error {
		for _, item := range req {
			branch, err := h.createBranchInternal(ctx, q, item)
			if err != nil {
				return err
			}
			branches = append(branches, branch)
		}
		return nil
	})

	if err != nil {
		if err.Error() == "department not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, branches)
}

func (h *branchHandler) createBranchInternal(ctx *gin.Context, q sqlc.Querier, req CreateBranchRequest) (sqlc.Branch, error) {
	// Better Approach: Associate by Department Name instead of requiring UUID from client
	dept, err := q.GetDepartmentByName(ctx, strings.ToLower(req.DepartmentName))
	if err != nil {
		return sqlc.Branch{}, fmt.Errorf("department not found")
	}

	arg := sqlc.CreateBranchParams{
		Name:         req.Name,
		Code:         strings.ToUpper(req.Code),
		DepartmentID: dept.ID,
	}

	return q.CreateBranch(ctx, arg)
}
