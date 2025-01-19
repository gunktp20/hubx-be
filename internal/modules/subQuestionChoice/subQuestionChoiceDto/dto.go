package dto

import (
	"time"

	"github.com/gunktp20/digital-hubx-be/pkg/models"
)

type (
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

	CreateSubQuestionChoicesRes struct {
		ID            string    `json:"id"`
		Description   string    `json:"description"`
		SubQuestionID string    `json:"sub_question_id"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
	}
)
