package dto

import (
	"time"

	"github.com/gunktp20/digital-hubx-be/pkg/models"
)

type (
	// CreateClassSessionReq represents the request body for creating a class session.
	// @Description Request body for creating a class session.
	// @Param class_id body string true "Class ID" format(uuid)
	// @Param date body string true "Session date" format(date-time)
	// @Param max_capacity body int true "Maximum capacity"
	// @Param start_time body string true "Start time" format(date-time)
	// @Param end_time body string true "End time" format(date-time)
	// @Param location body string true "Location"
	CreateClassSessionReq struct {
		ClassID     string    `json:"class_id" validate:"required,uuid"`
		Date        time.Time `json:"date" validate:"required"`
		MaxCapacity int       `json:"max_capacity" validate:"required"`
		// ClassSessionStatus models.ClassSessionStatus `json:"class_session_status" validate:"required"`
		StartTime time.Time `json:"start_time" validate:"required"`
		EndTime   time.Time `json:"end_time" validate:"required"`
		Location  string    `json:"location" validate:"required"`
	}

	// SetMaxCapacityReq represents the request body for updating max capacity.
	// @Description Request body for updating max capacity for a class session.
	// @Param new_capacity body int true "New max capacity"
	CreateClassSessionRes struct {
		ID                   string                    `json:"id"`
		ClassID              string                    `json:"class_id"`
		Date                 time.Time                 `json:"date" `
		MaxCapacity          int                       `json:"max_capacity" `
		CancellationDeadline time.Time                 `json:"cancellation_deadline"`
		ClassSessionStatus   models.ClassSessionStatus `json:"class_session_status" `
		StartTime            time.Time                 `json:"start_time" `
		EndTime              time.Time                 `json:"end_time" `
		Location             string                    `json:"location" `
		CreatedAt            time.Time                 `json:"created_at"`
		UpdatedAt            time.Time                 `json:"updated_at"`
	}

	// UpdateClassSessionLocation represents the request body for updating a class session's location.
	// @Description Request body for updating a class session's location.
	// @Param new_location body string true "New location"
	UpdateClassSessionLocation struct {
		NewLocation string `json:"new_location" validate:"required"`
	}

	// CreateClassSessionRes represents the response body after creating a class session.
	// @Description Response body for creating a class session.
	ClassSessionsRes struct {
		ID                 string                    `json:"id"`
		ClassID            string                    `json:"class_id"`
		Class              models.Class              `json:"class"`
		Date               time.Time                 `json:"date" `
		MaxCapacity        int                       `json:"max_capacity" `
		ClassSessionStatus models.ClassSessionStatus `json:"class_session_status" `
		StartTime          time.Time                 `json:"start_time" `
		EndTime            time.Time                 `json:"end_time" `
		Location           string                    `json:"location" `
		RemainingSeats     int                       `json:"remaining_seats" `
		CreatedAt          time.Time                 `json:"created_at"`
		UpdatedAt          time.Time                 `json:"updated_at"`
	}

	// SetMaxCapacityReq represents the request body for updating max capacity.
	// @Description Request body for updating max capacity for a class session.
	// @Param new_capacity body int true "New max capacity"
	SetMaxCapacityReq struct {
		NewCapacity int `json:"new_capacity" validate:"required"`
	}
)
