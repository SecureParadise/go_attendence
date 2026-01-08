package handlers

import (
	"encoding/csv"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/SecureParadise/go_attendence/internal/api/middleware"
	"github.com/SecureParadise/go_attendence/internal/db"
	"github.com/SecureParadise/go_attendence/internal/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type attendanceHandler struct {
	store db.Store
}

func NewAttendanceHandler(store db.Store) *attendanceHandler {
	return &attendanceHandler{store: store}
}

type MarkAttendanceRequest struct {
	StudentID  uuid.UUID             `json:"student_id" binding:"required"`
	SubjectID  uuid.UUID             `json:"subject_id" binding:"required"`
	TeacherID  uuid.UUID             `json:"teacher_id" binding:"required"`
	SemesterID uuid.UUID             `json:"semester_id" binding:"required"`
	Date       time.Time             `json:"date" binding:"required"`
	Status     sqlc.AttendanceStatus `json:"status" binding:"required"`
	Method     sqlc.AttendanceMethod `json:"method" binding:"required"`
	Remarks    string                `json:"remarks"`
}

func (h *attendanceHandler) MarkAttendance(ctx *gin.Context) {
	var req MarkAttendanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	arg := sqlc.CreateAttendanceParams{
		StudentID:  req.StudentID,
		SubjectID:  req.SubjectID,
		TeacherID:  req.TeacherID,
		SemesterID: req.SemesterID,
		Date: pgtype.Date{
			Time:  req.Date,
			Valid: true,
		},
		Status: req.Status,
		Method: req.Method,
		Remarks: pgtype.Text{
			String: req.Remarks,
			Valid:  req.Remarks != "",
		},
	}

	attendance, err := h.store.CreateAttendance(ctx, arg)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, attendance)
}

type CheckInRequest struct {
	StudentID  uuid.UUID             `json:"student_id" binding:"required"`
	SubjectID  uuid.UUID             `json:"subject_id" binding:"required"`
	TeacherID  uuid.UUID             `json:"teacher_id" binding:"required"`
	SemesterID uuid.UUID             `json:"semester_id" binding:"required"`
	Method     sqlc.AttendanceMethod `json:"method" binding:"required"`
}

func (h *attendanceHandler) CheckIn(ctx *gin.Context) {
	var req CheckInRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	now := time.Now()
	arg := sqlc.CreateAttendanceParams{
		StudentID:  req.StudentID,
		SubjectID:  req.SubjectID,
		TeacherID:  req.TeacherID,
		SemesterID: req.SemesterID,
		Date: pgtype.Date{
			Time:  now,
			Valid: true,
		},
		CheckIn: pgtype.Timestamptz{
			Time:  now,
			Valid: true,
		},
		Status: sqlc.AttendanceStatusPresent,
		Method: req.Method,
	}

	attendance, err := h.store.CreateAttendance(ctx, arg)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, attendance)
}

type CheckOutRequest struct {
	AttendanceID uuid.UUID `json:"attendance_id" binding:"required"`
}

func (h *attendanceHandler) CheckOut(ctx *gin.Context) {
	var req CheckOutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	now := time.Now()
	arg := sqlc.UpdateAttendanceParams{
		ID: req.AttendanceID,
		CheckOut: pgtype.Timestamptz{
			Time:  now,
			Valid: true,
		},
	}

	attendance, err := h.store.UpdateAttendance(ctx, arg)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, attendance)
}

type GetReportRequest struct {
	SemesterID uuid.UUID `form:"semester_id" binding:"required"`
	StartDate  time.Time `form:"start_date" binding:"required" time_format:"2006-01-02"`
	EndDate    time.Time `form:"end_date" binding:"required" time_format:"2006-01-02"`
}

func (h *attendanceHandler) GetAttendanceReport(ctx *gin.Context) {
	var req GetReportRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(err)
		return
	}

	arg := sqlc.ListAttendanceForReportParams{
		SemesterID: req.SemesterID,
		Date: pgtype.Date{
			Time:  req.StartDate,
			Valid: true,
		},
		Date_2: pgtype.Date{
			Time:  req.EndDate,
			Valid: true,
		},
	}

	report, err := h.store.ListAttendanceForReport(ctx, arg)
	if err != nil {
		ctx.Error(err)
		return
	}

	format := ctx.Query("format")
	if format == "csv" {
		ctx.Header("Content-Disposition", "attachment; filename=attendance_report.csv")
		ctx.Header("Content-Type", "text/csv")

		writer := csv.NewWriter(ctx.Writer)
		defer writer.Flush()

		writer.Write([]string{"Date", "Roll No", "Student Name", "Subject", "Teacher", "Status", "Check-in", "Check-out", "Method"})

		for _, row := range report {
			writer.Write([]string{
				row.Date.Time.Format("2006-01-02"),
				row.RollNo,
				fmt.Sprintf("%s %s", row.FirstName, row.LastName),
				row.SubjectName,
				fmt.Sprintf("%s %s", row.TeacherFirstName, row.TeacherLastName),
				string(row.Status),
				row.CheckIn.Time.Format("15:04:05"),
				row.CheckOut.Time.Format("15:04:05"),
				string(row.Method),
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, report)
}

func (h *attendanceHandler) DeviceMarkAttendance(ctx *gin.Context) {
	type DeviceRequest struct {
		StudentID uuid.UUID             `json:"student_id" binding:"required"`
		Method    sqlc.AttendanceMethod `json:"method" binding:"required"`
	}

	var req DeviceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	// 1. Find the active class session
	// This query needs to be defined in core_attendance.sql
	// For simplicity, let's assume we search for ANY active session now
	// In a real scenario, the device might be associated with a room/subject.
	// Here we just find an active session the student should be in.

	// Find session by student's current enrollment
	// This is a bit complex for a single query, so let's use a simpler approach for now
	// and assume there's only one active session globally or per room.

	// For the sake of the requirement: "find the active class_session"
	// We'll search for a session that started within the last 90 minutes.
	session, err := h.store.GetActiveSessionForStudent(ctx, req.StudentID)
	if err != nil {
		ctx.Error(middleware.NewAPIError(http.StatusNotFound, "no active class session found for student", err))
		return
	}

	now := time.Now()
	diff := now.Sub(session.ActualStart)
	minutes := diff.Minutes()

	var score float64
	status := sqlc.AttendanceStatusPresent

	if minutes <= 15 {
		score = 1.0
	} else if minutes <= 40 {
		score = 0.8
		status = sqlc.AttendanceStatusLate
	} else if minutes <= 90 {
		score = 0.6
		status = sqlc.AttendanceStatusLate
	} else {
		score = 0.0
		status = sqlc.AttendanceStatusAbsent
	}

	arg := sqlc.CreateAttendanceRecordParams{
		StudentID: req.StudentID,
		SessionID: session.ID,
		ScanTime: pgtype.Timestamptz{
			Time:  now,
			Valid: true,
		},
		Score:  pgtype.Numeric{Int: big.NewInt(int64(score * 100)), Exp: -2, Valid: true},
		Status: status,
		Method: req.Method,
	}

	record, err := h.store.CreateAttendanceRecord(ctx, arg)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, record)
}

func (h *attendanceHandler) GetStudentPercentage(ctx *gin.Context) {
	studentID, err := uuid.Parse(ctx.Param("student_id"))
	if err != nil {
		ctx.Error(middleware.NewAPIError(http.StatusBadRequest, "invalid student id", err))
		return
	}

	semesterID, err := uuid.Parse(ctx.Query("semester_id"))
	if err != nil {
		ctx.Error(middleware.NewAPIError(http.StatusBadRequest, "invalid semester id", err))
		return
	}

	arg := sqlc.GetStudentAttendancePercentageParams{
		StudentID:  studentID,
		SemesterID: semesterID,
	}

	percentage, err := h.store.GetStudentAttendancePercentage(ctx, arg)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, percentage)
}
