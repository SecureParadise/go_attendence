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
	studentHandler := handlers.NewStudentHandler(store)
	teacherHandler := handlers.NewTeacherHandler(store)

	// Create a single user
	router.POST("/register", userHandler.CreateUser)
	// User login
	router.POST("/login", userHandler.Login)
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
	// Create a single student
	router.POST("/student_reg", studentHandler.CreateStudent)
	// Get student by roll number
	router.GET("/student/:roll_no", studentHandler.GetStudentByRollNo)
	// Create a single teacher
	router.POST("/teacher_reg", teacherHandler.CreateTeacher)
	// Get teacher by card number
	router.GET("/teacher/:card_no", teacherHandler.GetTeacherByCardNo)
}
