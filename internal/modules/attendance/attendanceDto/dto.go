package dto

import (
	"github.com/gunktp20/digital-hubx-be/pkg/models"
)

type (
	CreateAttendanceReq struct {
		UserEmail      string `json:"user_email" validate:"required,email"`
		ClassID        string `json:"class_id"`
		ClassSessionID string `json:"class_session_id" validate:"required,uuid"`
	}

	AttendanceClass struct {
		Title       string           `json:"title"`
		Description string           `json:"description"`
		ClassTier   models.ClassTier `json:"class_tier"`
		ClassLevel  int              `json:"class_level"`
	}

	CreateAttendanceRes struct {
		ID             string `json:"id"`
		UserEmail      string `json:"user_email"`
		ClassID        string `json:"class_id" validate:"required,uuid"`
		ClassSessionID string `json:"class_session_id" validate:"required,uuid"`
	}
)
