package handlers

import (
	"net/http"
	"strings"

	"github.com/SecureParadise/go_attendence/internal/db"
	"github.com/SecureParadise/go_attendence/internal/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateDepartmentRequest struct {
	Name     string      `json:"name" binding:"required"`
	HodName  pgtype.Text `json:"hod_name"`
	DhodName pgtype.Text `json:"dhod_name"`
}

type departmentHandler struct {
	store db.Store
}

func NewDepartmentHandler(store db.Store) *departmentHandler {
	return &departmentHandler{store: store}
}

func (h *departmentHandler) CreateDepartment(ctx *gin.Context) {
	var req CreateDepartmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	department, err := h.createDepartmentInternal(ctx, h.store, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, department)
}

func (h *departmentHandler) BulkCreateDepartments(ctx *gin.Context) {
	var req []CreateDepartmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var departments []sqlc.Department
	err := h.store.WithTx(ctx, func(q *sqlc.Queries) error {
		for _, item := range req {
			dept, err := h.createDepartmentInternal(ctx, q, item)
			if err != nil {
				return err
			}
			departments = append(departments, dept)
		}
		return nil
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, departments)
}

func (h *departmentHandler) createDepartmentInternal(ctx *gin.Context, q sqlc.Querier, req CreateDepartmentRequest) (sqlc.Department, error) {
	arg := sqlc.CreateDepartmentParams{
		Name: strings.ToLower(req.Name),
	}

	if req.HodName.Valid {
		arg.HodName = pgtype.Text{
			String: strings.ToLower(req.HodName.String),
			Valid:  true,
		}
	}

	if req.DhodName.Valid {
		arg.DhodName = pgtype.Text{
			String: strings.ToLower(req.DhodName.String),
			Valid:  true,
		}
	}

	return q.CreateDepartment(ctx, arg)
}
