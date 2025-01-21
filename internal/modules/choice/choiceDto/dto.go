package dto

import (
	"github.com/gunktp20/digital-hubx-be/pkg/models"
)

type (
	// CreateChoiceReq represents the request body for creating a choice.
	// @Description Request body for creating a new choice.
	// @Param description query string true "Description of the choice"
	// @Param question_id query string true "ID of the question this choice belongs to" format(uuid)
	CreateChoiceReq struct {
		Description string `json:"description" validate:"required"`      // Choice description
		QuestionID  string `json:"question_id" validate:"required,uuid"` // Question ID
	}

	// ChoiceClass represents the details of a class related to a choice.
	// @Description Class details related to a choice.
	// @Param title query string true "Class title"
	// @Param description query string true "Class description"
	// @Param class_tier query string true "Class tier" enum(Essential,Literacy,Mastery)
	// @Param class_level query int true "Class level"
	ChoiceClass struct {
		Title       string           `json:"title"`       // Class title
		Description string           `json:"description"` // Class description
		ClassTier   models.ClassTier `json:"class_tier"`  // Class tier (e.g., Essential, Literacy, Mastery)
		ClassLevel  int              `json:"class_level"` // Class level
	}

	// CreateChoiceRes represents the response body after creating a choice.
	// @Description Response body after creating a new choice.
	// @Param id query string true "Choice ID" format(uuid)
	// @Param description query string true "Description of the choice"
	// @Param question_id query string true "ID of the question this choice belongs to" format(uuid)
	CreateChoiceRes struct {
		ID          string `json:"id"`          // Choice ID
		Description string `json:"description"` // Choice description
		QuestionID  string `json:"question_id"` // Question ID
	}
)
