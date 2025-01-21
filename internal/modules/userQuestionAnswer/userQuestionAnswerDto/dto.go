package dto

import (
	"time"

	"github.com/gunktp20/digital-hubx-be/pkg/models"
)

type (

	// SubQuestionsAnswer represents an answer for a sub-question.
	// @Description Details of a sub-question answer.
	// @Param sub_question_id body string true "Sub-question ID"
	// @Param answer_text body string false "Answer text"
	// @Param selected_sub_question_choice_id body string false "Selected sub-question choice ID"
	SubQuestionsAnswer struct {
		SubQuestionID  string `json:"sub_question_id"`
		Description    string `json:"description"`
		ParentChoiceID string `json:"parent_choice_id"`
		QuestionType   string `json:"question_type"`
		AnswerText     string `json:"answer_text"`

		SelectedSubQuestionChoiceID          string `json:"selected_sub_question_choice_id"`
		SelectedSubQuestionChoiceDescription string `json:"selected_sub_question_choice_description"`
	}

	CreateUserQuestionAnswerReq struct {
		QuestionID       string `json:"question_id" validate:"required"`
		SelectedChoiceID string `json:"selected_choice_id"`
		AnswerText       string `json:"answer_text"`

		SubQuestionsAnswers []SubQuestionsAnswer `json:"sub_question_answers"`

		// ? user don't need to need to provide
	}

	UserQuestionAnswerClass struct {
		Title       string           `json:"title"`
		Description string           `json:"description"`
		ClassTier   models.ClassTier `json:"class_tier"`
		ClassLevel  int              `json:"class_level"`
	}

	CreateUserQuestionAnswerRes struct {
		ID         string  `json:"id"`
		QuestionID string  `json:"question_id"`
		ChoiceID   *string `json:"choice_id"`
		ClassID    string  `json:"class_id"`
		AnswerText string  `json:"answer_text"`
		Email      string  `json:"email"`

		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	SubQChoiceRes struct {
		ID          string `json:"id"`
		Description string `json:"description"`
	}

	SubQuestionAnswerRes struct {
		SubQuestionID                        string               `json:"sub_question_id"`
		SubQuestionDescription               string               `json:"sub_question_description"`
		QuestionType                         models.QuestionTypes `json:"question_type"`
		SelectedSubQuestionChoiceID          string               `json:"selected_sub_question_choice_id"`
		SelectedSubQuestionChoiceDescription string               `json:"selected_sub_question_choice_description"`
		AnswerText                           string               `json:"answer_text"`
	}

	GetUserQuestionAnswersChoice struct {
		ID                 string                 `json:"id"`
		Description        string                 `json:"description"`
		SubQuestionAnswers []SubQuestionAnswerRes `json:"sub_question_answers"`
	}

	GetUserQuestionAnswersQuestion struct {
		Description string `json:"description"`
	}

	// GetUserQuestionAnswerRes represents the response body for a user question answer.
	// @Description Response body for a user question answer.
	// @Param question_id body string true "Question ID"
	// @Param selected_choice_id body string false "Selected choice ID"
	// @Param answer_text body string false "Answer text"
	GetUserQuestionAnswerRes struct {
		ID               string                         `json:"id"`
		QuestionID       string                         `json:"question_id"`
		Question         GetUserQuestionAnswersQuestion `json:"question"`
		SelectedChoiceID string                         `json:"selected_choice_id"`
		SelectedChoice   GetUserQuestionAnswersChoice   `json:"selected_choice"`
		ClassID          string                         `json:"class_id"`
		AnswerText       string                         `json:"answer_text"`
		QuestionType     models.QuestionTypes           `json:"question_type"`
		Email            string                         `json:"email"`

		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)
