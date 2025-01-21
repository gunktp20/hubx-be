package dto

import (
	"time"

	"github.com/gunktp20/digital-hubx-be/pkg/models"
)

type (

	// CreateSubQuestionChoicesReq represents the request body for creating a sub-question choice.
	// @Description Request body for creating a sub-question choice.
	// @Param description body string true "Description of the sub-question choice"
	// @Param sub_question_id body string true "Sub-question ID" format(uuid)
	CreateSubQuestionChoicesReq struct {
		Description   string `json:"description" validate:"required"`
		SubQuestionID string `json:"sub_question_id" validate:"required,uuid"`
	}

	SubQuestionChoicesClass struct {
		Title       string           `json:"title"`
		Description string           `json:"description"`
		ClassTier   models.ClassTier `json:"class_tier"`
		ClassLevel  int              `json:"class_level"`
	}

	// CreateSubQuestionChoicesRes represents the response body after creating a sub-question choice.
	// @Description Response body for creating a sub-question choice.
	// @Param id body string true "Sub-question choice ID" format(uuid)
	// @Param description body string true "Description of the sub-question choice"
	// @Param sub_question_id body string true "Sub-question ID" format(uuid)
	// @Param created_at body string true "Timestamp of creation" format(date-time)
	// @Param updated_at body string true "Timestamp of last update" format(date-time)
	CreateSubQuestionChoicesRes struct {
		ID            string    `json:"id"`
		Description   string    `json:"description"`
		SubQuestionID string    `json:"sub_question_id"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
	}
)
