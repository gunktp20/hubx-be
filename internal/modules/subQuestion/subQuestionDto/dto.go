package dto

import (
	"time"

	"github.com/gunktp20/digital-hubx-be/pkg/models"
)

type (
	CreateSubQuestionReq struct {
		Description  string               `json:"description" validate:"required"`
		ChoiceID     string               `json:"choice_id" validate:"required,uuid"`
		QuestionType models.QuestionTypes `json:"question_type" validate:"required"`

		// ? User can ignore these
		QuestionID string `json:"question_id"`

		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	SubQuestionClass struct {
		Title       string           `json:"title"`
		Description string           `json:"description"`
		ClassTier   models.ClassTier `json:"class_tier"`
		ClassLevel  int              `json:"class_level"`
	}

	QuestionRes struct {
		Description string `json:"description"`
	}

	ChoiceRes struct {
		Description string `json:"description"`
	}

	CreateSubQuestionChoiceRes struct {
		Description string `json:"description"`
	}

	CreateSubQuestionRes struct {
		ID          string      `json:"id"`
		Description string      `json:"description"`
		QuestionID  string      `json:"question_id"`
		Question    QuestionRes `json:"question"`
		ChoiceID    string      `json:"choice_id"`
		Choice      ChoiceRes   `json:"choice"`

		QuestionType models.QuestionTypes `json:"question_type"`
		CreatedAt    time.Time            `json:"created_at"`
		UpdatedAt    time.Time            `json:"updated_at"`

		SubQuestionChoices []CreateSubQuestionChoiceRes `json:"sub_question_choices"`
	}
)
