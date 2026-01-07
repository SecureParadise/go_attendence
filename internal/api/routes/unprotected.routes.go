package routes

import (
	"github.com/SecureParadise/go_attendence/internal/api/handlers"
	"github.com/SecureParadise/go_attendence/internal/db"
	"github.com/gin-gonic/gin"
)

func SetupUnProtectedRoutes(router *gin.Engine, store db.Store) {
	userHandler := handlers.NewUserHandler(store)
	deptHandler := handlers.NewDepartmentHandler(store)
	branchHandler := handlers.NewBranchHandler(store)
	semesterHandler := handlers.NewSemesterHandler(store)

	// Create a single user
	router.POST("/register", userHandler.CreateUser)
	// Create a single department
	router.POST("/dept_reg", deptHandler.CreateDepartment)
	// Create multiple departments
	router.POST("/dept_bulk_reg", deptHandler.BulkCreateDepartments)
	// Create a single branch
	router.POST("/branch_reg", branchHandler.CreateBranch)
	// Create multiple branches
	router.POST("/branch_bulk_reg", branchHandler.BulkCreateBranches)
	// Create a single semester
	router.POST("/semester_reg", semesterHandler.CreateSemester)
}
