package dto

import (
	"time"

	"github.com/gunktp20/digital-hubx-be/pkg/models"
)

type (
	// CreateQuestionReq represents the request body for creating a question.
	// @Description Request body for creating a question.
	// @Param description body string true "Question description"
	// @Param class_id body string true "Class ID" format(uuid)
	// @Param question_type body string true "Question type" enum(SingleChoice,MultipleChoice)
	CreateQuestionReq struct {
		Description  string               `json:"description" validate:"required"`
		ClassID      string               `json:"class_id" validate:"required,uuid"`
		QuestionType models.QuestionTypes `json:"question_type" validate:"required"`
	}

	QuestionClass struct {
		Title       string           `json:"title"`
		Description string           `json:"description"`
		ClassTier   models.ClassTier `json:"class_tier"`
		ClassLevel  int              `json:"class_level"`
	}

	// CreateQuestionRes represents the response body after creating a question.
	// @Description Response body for creating a question.
	// @Param id body string true "Question ID" format(uuid)
	// @Param description body string true "Question description"
	// @Param class_id body string true "Class ID" format(uuid)
	// @Param question_type body string true "Question type" enum(SingleChoice,MultipleChoice)
	// @Param created_at body string true "Timestamp of creation" format(date-time)
	// @Param updated_at body string true "Timestamp of last update" format(date-time)
	// @Param choices body []models.Choice true "List of associated choices"
	CreateQuestionRes struct {
		ID           string               `json:"id"`
		Description  string               `json:"description"`
		ClassID      string               `json:"class_id"`
		Class        QuestionClass        `json:"class"`
		QuestionType models.QuestionTypes `json:"question_type"`
		CreatedAt    time.Time            `json:"created_at"`
		UpdatedAt    time.Time            `json:"updated_at"`

		Choices []models.Choice `json:"choices" validate:"required"`
	}
)
