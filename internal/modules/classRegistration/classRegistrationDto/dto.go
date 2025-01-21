package dto

import (
	"time"

	"github.com/gunktp20/digital-hubx-be/pkg/models"
)

type (

	// CreateClassRegistrationReq represents the request body for creating a class registration.
	// @Description Request body for creating a class registration.
	// @Param class_id body string true "Class ID" format(uuid)
	// @Param class_session_id body string true "Class session ID" format(uuid)
	CreateClassRegistrationReq struct {
		ClassID        string `json:"class_id" validate:"required,uuid"`
		ClassSessionID string `json:"class_session_id" validate:"required,uuid"`
	}

	// CreateClassRegistrationRes represents the response body after creating a class registration.
	// @Description Response body for a successful class registration creation.
	CreateClassRegistrationRes struct {
		ID                   string           `json:"id"`
		Email                string           `json:"email"`
		ClassID              string           `json:"class_id"`
		ClassSessionID       string           `json:"class_session_id"`
		UnattendedQuota      int              `json:"unattended_quota"`
		IsBanned             bool             `json:"is_banned"`
		CancellationDeadline time.Time        `json:"cancellation_deadline"`
		RegisteredAt         time.Time        `json:"registered_at"`
		RegStatus            models.RegStatus `json:"reg_status"`
		CreatedAt            time.Time        `json:"created_at"`
		UpdatedAt            time.Time        `json:"updated_at"`
	}

	GetUserRegistrationsRes struct {
		ID                   string           `json:"id"`
		Email                string           `json:"email"`
		ClassID              string           `json:"class_id"`
		ClassSessionID       string           `json:"class_session_id"`
		UnattendedQuota      int              `json:"unattended_quota"`
		IsBanned             bool             `json:"is_banned"`
		CancellationDeadline time.Time        `json:"cancellation_deadline"`
		RegisteredAt         time.Time        `json:"registered_at"`
		RegStatus            models.RegStatus `json:"reg_status"`
		CreatedAt            time.Time        `json:"created_at"`
		UpdatedAt            time.Time        `json:"updated_at"`
	}

	// ResetCancelledQuotaReq represents the request body for resetting a user's cancellation quota.
	// @Description Request body for resetting a user's cancellation quota.
	// @Param user_email body string true "User's email" format(email)
	// @Param class_id body string true "Class ID" format(uuid)
	ResetCancelledQuotaReq struct {
		UserEmail string `json:"user_email" validate:"required,email"`
		ClassID   string `json:"class_id" validate:"required,uuid"`
	}
)
