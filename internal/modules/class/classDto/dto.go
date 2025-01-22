package dto

import (
	"time"

	"github.com/gunktp20/digital-hubx-be/pkg/models"
)

type (
	// CreateClassReq represents the request body for creating a class.
	// @Description Request body for creating a class.
	// @Param title query string true "Title of the class" minlength(5)
	// @Param description query string true "Description of the class" minlength(5)
	// @Param cover_image query string false "Cover image URL of the class"
	// @Param class_category_id query string true "ID of the class category"
	// @Param class_level query int false "Level of the class"
	// @Param class_tier query string true "Tier of the class" enum(Essential,Literacy,Mastery)
	CreateClassReq struct {
		Title           string           `json:"title" validate:"required,min=5"`
		Description     string           `json:"description" validate:"required,min=5"`
		CoverImage      string           `json:"cover_image"`
		ClassCategoryID string           `json:"class_category_id" validate:"required"`
		ClassLevel      int              `json:"class_level"`
		ClassTier       models.ClassTier `json:"class_tier" validate:"required"`
	}

	// CreateClassRes represents the response body after creating a class.
	// @Description Response body after creating a class.
	// @Param id query string true "Unique identifier for the class" format(uuid)
	// @Param title query string true "Title of the class"
	// @Param description query string true "Description of the class"
	// @Param cover_image query string false "Cover image URL of the class"
	// @Param class_category_id query string true "ID of the class category"
	// @Param class_level query int false "Level of the class"
	// @Param class_tier query string true "Tier of the class" enum(Essential,Literacy,Mastery)
	// @Param is_active query boolean true "Is the class currently active?"
	// @Param is_remove query boolean true "Is the class marked as removed?"
	// @Param created_at query string true "Timestamp of when the class was created" format(date-time)
	// @Param updated_at query string true "Timestamp of when the class was last updated" format(date-time)
	CreateClassRes struct {
		ID              string           `json:"id"`
		Title           string           `json:"title"`
		Description     string           `json:"description"`
		CoverImage      string           `json:"cover_image"`
		ClassCategoryID string           `json:"class_category_id"`
		ClassLevel      int              `json:"class_level"`
		ClassTier       models.ClassTier `json:"class_tier"`
		IsActive        bool             `json:"is_active"`
		IsRemove        bool             `json:"is_remove"`
		CreatedAt       time.Time        `json:"created_at"`
		UpdatedAt       time.Time        `json:"updated_at"`
	}

	// UpdateClassReq represents the request body for updating class details.
	// @Description Request body for updating class details.
	// @Param title query string false "New title for the class"
	// @Param description query string false "New description for the class"
	// @Param class_category_name query string false "New name of the class category"
	// @Param class_tier query string false "New tier of the class" enum(Essential,Literacy,Mastery)
	// @Param class_level query int false "New level of the class"
	UpdateClassReq struct {
		Title             *string `json:"title"`
		Description       *string `json:"description"`
		ClassCategoryName *string `json:"class_category_name"`
		ClassTier         *string `json:"class_tier"`
		ClassLevel        *int    `json:"class_level"`
	}

	ClassRes struct {
		ID             string                `json:"id"`
		Title          string                `json:"title"`
		Description    string                `json:"description"`
		CoverImage     string                `json:"cover_image"`
		ClassTier      models.ClassTier      `json:"class_tier"`
		ClassLevel     int                   `json:"class_level"`
		IsActive       bool                  `json:"is_active"`
		IsRemove       bool                  `json:"is_remove"`
		EnableQuestion bool                  `json:"enable_question"`
		Order          int                   `json:"order"`
		CreatedAt      time.Time             `json:"created_at"`
		UpdatedAt      time.Time             `json:"updated_at"`
		ClassCategory  models.ClassCategory  `json:"class_category"`
		ClassSessions  []models.ClassSession `json:"class_sessions"`
		IsRegistered   bool                  `json:"is_registered"` // Registration status
	}
)
