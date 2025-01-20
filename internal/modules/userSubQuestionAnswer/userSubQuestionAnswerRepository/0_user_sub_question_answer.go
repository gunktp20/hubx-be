package repository

import (
	userSubQuestionAnswerDto "github.com/gunktp20/digital-hubx-be/internal/modules/userSubQuestionAnswer/userSubQuestionAnswerDto"
	"github.com/gunktp20/digital-hubx-be/pkg/models"
	"gorm.io/gorm"
)

type (
	UserSubQuestionAnswerRepositoryService interface {
		CreateUserSubQuestionAnswer(tx *gorm.DB, createUserSubQuestionAnswerReq *userSubQuestionAnswerDto.CreateUserSubQuestionAnswerReq, email string) (*userSubQuestionAnswerDto.CreateUserSubQuestionAnswerRes, error)
	}
)

func (r *userSubQuestionAnswerGormRepository) CreateUserSubQuestionAnswer(tx *gorm.DB, createUserSubQuestionAnswerReq *userSubQuestionAnswerDto.CreateUserSubQuestionAnswerReq, email string) (*userSubQuestionAnswerDto.CreateUserSubQuestionAnswerRes, error) {

	var subQuestionChoiceID *string
	if createUserSubQuestionAnswerReq.SubQuestionChoiceID != "" {
		subQuestionChoiceID = &createUserSubQuestionAnswerReq.SubQuestionChoiceID
	}

	// แปลง AnswerText เป็น *string
	var answerText *string
	if createUserSubQuestionAnswerReq.AnswerText != "" {
		answerText = &createUserSubQuestionAnswerReq.AnswerText
	}

	userSubQuestionAnswer := models.UserSubQuestionAnswer{
		UserEmail:           email,
		SubQuestionID:       createUserSubQuestionAnswerReq.SubQuestionID,
		SubQuestionChoiceID: subQuestionChoiceID,
		ClassID:             createUserSubQuestionAnswerReq.ClassID,
		AnswerText:          answerText, // ใช้ค่า pointer
	}

	if err := tx.Create(&userSubQuestionAnswer).Error; err != nil {
		return &userSubQuestionAnswerDto.CreateUserSubQuestionAnswerRes{}, err
	}

	return &userSubQuestionAnswerDto.CreateUserSubQuestionAnswerRes{
		ID:                  userSubQuestionAnswer.ID,
		Email:               email,
		SubQuestionID:       userSubQuestionAnswer.SubQuestionID,
		SubQuestionChoiceID: subQuestionChoiceID,
		ClassID:             userSubQuestionAnswer.ClassID,
		AnswerText:          answerText,
		CreatedAt:           userSubQuestionAnswer.CreatedAt,
		UpdatedAt:           userSubQuestionAnswer.UpdatedAt,
	}, nil
}
