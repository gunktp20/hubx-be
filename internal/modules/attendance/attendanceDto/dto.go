package dto

import (
	"github.com/gunktp20/digital-hubx-be/pkg/models"
)

// CreateAttendanceReq represents the request body for creating an attendance record.
// @Description Request body for creating an attendance record.
// @Param user_email query string true "User's email address" format(email)
// @Param class_id query string false "ID of the class"
// @Param class_session_id query string true "ID of the class session" format(uuid)
type CreateAttendanceReq struct {
	UserEmail      string `json:"user_email" validate:"required,email"`      // User's email
	ClassID        string `json:"class_id"`                                  // Class ID
	ClassSessionID string `json:"class_session_id" validate:"required,uuid"` // Class session ID
}

// AttendanceClass represents a class's details for attendance.
// @Description Details of a class for attendance purposes.
// @Param title query string true "Class title"
// @Param description query string true "Class description"
// @Param class_tier query string true "Class tier" enum(Essential,Literacy,Mastery)
// @Param class_level query int true "Class level"
type AttendanceClass struct {
	Title       string           `json:"title"`       // Class title
	Description string           `json:"description"` // Class description
	ClassTier   models.ClassTier `json:"class_tier"`  // Class tier (e.g., Essential, Literacy, Mastery)
	ClassLevel  int              `json:"class_level"` // Class level
}

// CreateAttendanceRes represents the response body after creating an attendance record.
// @Description Response body after creating an attendance record.
// @Param id query string true "Attendance ID" format(uuid)
// @Param user_email query string true "User's email address" format(email)
// @Param class_id query string true "ID of the class" format(uuid)
// @Param class_session_id query string true "ID of the class session" format(uuid)
type CreateAttendanceRes struct {
	ID             string `json:"id"`               // Attendance ID
	UserEmail      string `json:"user_email"`       // User's email
	ClassID        string `json:"class_id"`         // Class ID
	ClassSessionID string `json:"class_session_id"` // Class session ID
}
