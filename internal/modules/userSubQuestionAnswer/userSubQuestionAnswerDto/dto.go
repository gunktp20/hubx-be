package dto

import (
	"time"

	"github.com/gunktp20/digital-hubx-be/pkg/models"
)

type (
	CreateUserSubQuestionAnswerReq struct {
		SubQuestionChoiceID string `json:"sub_question_choice_id" validate:"required"`
		ClassID             string `json:"class_id"`
		AnswerText          string `json:"answer_text"`

		// ? don't need to provide
		SubQuestionID string `json:"sub_question_id"`
	}

	UserSubQuestionAnswerClass struct {
		Title       string           `json:"title"`
		Description string           `json:"description"`
		ClassTier   models.ClassTier `json:"class_tier"`
		ClassLevel  int              `json:"class_level"`
	}

	CreateUserSubQuestionAnswerRes struct {
		ID                  string    `json:"id"`
		Email               string    `json:"email"`
		SubQuestionID       string    `json:"sub_question_id"`
		SubQuestionChoiceID *string   `json:"sub_question_choice_id"`
		ClassID             string    `json:"class_id"`
		AnswerText          string    `json:"answer_text"`
		CreatedAt           time.Time `json:"created_at"`
		UpdatedAt           time.Time `json:"updated_at"`
	}

	GetUserSubQuestionAnswersChoice struct {
		Description string `json:"description"`
	}

	GetUserSubQuestionAnswersQuestion struct {
		Description string `json:"description"`
	}

	GetUserSubQuestionAnswerRes struct {
		ID         string                            `json:"id"`
		QuestionID string                            `json:"question_id"`
		Question   GetUserSubQuestionAnswersQuestion `json:"question"`
		ChoiceID   string                            `json:"choice_id"`
		Choice     GetUserSubQuestionAnswersChoice   `json:"choice"`
		ClassID    string                            `json:"class_id"`
		AnswerText string                            `json:"answer_text"`
		Email      string                            `json:"email"`

		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)
