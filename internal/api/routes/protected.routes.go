// internal/api/routes/protected.routes.go
package routes

import (
	"github.com/SecureParadise/go_attendence/internal/api/handlers"
	"github.com/SecureParadise/go_attendence/internal/api/middleware"
	"github.com/SecureParadise/go_attendence/internal/auth"
	"github.com/SecureParadise/go_attendence/internal/config"
	"github.com/SecureParadise/go_attendence/internal/db"
	"github.com/SecureParadise/go_attendence/internal/db/sqlc"
	"github.com/gin-gonic/gin"
)

func SetupProtectedRoutes(router *gin.Engine, store db.Store, tokenMaker auth.Maker, config config.Config) {
	authRoutes := router.Group("/")
	authRoutes.Use(middleware.AuthMiddleware(tokenMaker))

	attendanceHandler := handlers.NewAttendanceHandler(store)
	userHandler := handlers.NewUserHandler(store, tokenMaker, config)
	studentHandler := handlers.NewStudentHandler(store)
	teacherHandler := handlers.NewTeacherHandler(store)

	// Admin only routes
	adminRoutes := authRoutes.Group("/").Use(middleware.RoleMiddleware(string(sqlc.UserroleAdmin)))
	adminRoutes.POST("/branch_reg", handlers.NewBranchHandler(store).CreateBranch)
	adminRoutes.POST("/branch_bulk_reg", handlers.NewBranchHandler(store).BulkCreateBranches)
	adminRoutes.POST("/dept_reg", handlers.NewDepartmentHandler(store).CreateDepartment)
	adminRoutes.POST("/dept_bulk_reg", handlers.NewDepartmentHandler(store).BulkCreateDepartments)
	adminRoutes.POST("/semester_reg", handlers.NewSemesterHandler(store).CreateSemester)

	// Teacher or Admin routes
	teacherAdminRoutes := authRoutes.Group("/").Use(middleware.RoleMiddleware(string(sqlc.UserroleTeacher), string(sqlc.UserroleAdmin)))
	teacherAdminRoutes.POST("/attendance/mark", attendanceHandler.MarkAttendance)
	teacherAdminRoutes.GET("/attendance/report", attendanceHandler.GetAttendanceReport)

	// Registration Completion (Protected by Auth, but specific to role)
	authRoutes.POST("/student_reg", handlers.NewStudentHandler(store).CreateStudent)
	authRoutes.POST("/teacher_reg", handlers.NewTeacherHandler(store).CreateTeacher)

	// Device endpoint for RFID/Fingerprint (high performance)
	authRoutes.POST("/attendance/device", attendanceHandler.DeviceMarkAttendance)
	// Student percentage
	authRoutes.GET("/attendance/student/:student_id/percentage", attendanceHandler.GetStudentPercentage)
	authRoutes.GET("/user/me", userHandler.GetUserMe)

	// Get student by roll number
	authRoutes.GET("/student/:roll_no", studentHandler.GetStudentByRollNo)
	// Get teacher by card number
	authRoutes.GET("/teacher/:card_no", teacherHandler.GetTeacherByCardNo)
}
