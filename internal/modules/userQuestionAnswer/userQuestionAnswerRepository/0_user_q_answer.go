package repository

import (
	"errors"

	userQuestionAnswerDto "github.com/gunktp20/digital-hubx-be/internal/modules/userQuestionAnswer/userQuestionAnswerDto"
	"github.com/gunktp20/digital-hubx-be/pkg/models"
	"gorm.io/gorm"
)

type (
	UserQuestionAnswerRepositoryService interface {
		CreateUserQuestionAnswer(tx *gorm.DB, createUserQuestionAnswerReq *models.UserQuestionAnswer) (*userQuestionAnswerDto.CreateUserQuestionAnswerRes, error)
		GetUserQuestionAnswersWithClassId(email, classID string, page int, limit int) (*[]userQuestionAnswerDto.GetUserQuestionAnswerRes, int64, error)
		GetUserQuestionAnswerById(userQuestionAnswerID string) (*models.UserQuestionAnswer, error)
		IsUserAnsweredThisQuestion(email, questionID string) (bool, error)
	}
)

func (r *userQuestionAnswerGormRepository) CreateUserQuestionAnswer(tx *gorm.DB, createUserQuestionAnswerReq *models.UserQuestionAnswer) (*userQuestionAnswerDto.CreateUserQuestionAnswerRes, error) {

	userQuestionAnswer := createUserQuestionAnswerReq

	if err := tx.Create(&userQuestionAnswer).Error; err != nil {
		return &userQuestionAnswerDto.CreateUserQuestionAnswerRes{}, err
	}

	return &userQuestionAnswerDto.CreateUserQuestionAnswerRes{
		ID:         userQuestionAnswer.ID,
		QuestionID: userQuestionAnswer.QuestionID,
		ChoiceID:   userQuestionAnswer.ChoiceID,
		// ClassID:    createUserQuestionAnswerReq.ClassID,
		AnswerText: userQuestionAnswer.AnswerText,
		Email:      userQuestionAnswer.UserEmail,
		CreatedAt:  userQuestionAnswer.CreatedAt,
		UpdatedAt:  userQuestionAnswer.UpdatedAt,
	}, nil
}

func (r *userQuestionAnswerGormRepository) GetUserQuestionAnswersWithClassId(email, classID string, page int, limit int) (*[]userQuestionAnswerDto.GetUserQuestionAnswerRes, int64, error) {
	var userQuestionAnswers []models.UserQuestionAnswer
	var total int64

	query := r.db.Model(&models.UserQuestionAnswer{}).Where("email = ? AND class_id = ?", email, classID).Count(&total)

	offset := (page - 1) * limit
	result := query.
		Preload("Choice").
		Preload("Question").
		Preload("Choice.SubQuestions.SubQuestionChoices").
		Limit(limit).
		Offset(offset).
		Find(&userQuestionAnswers)

	var userQuestionAnswerRes []userQuestionAnswerDto.GetUserQuestionAnswerRes

	for i := range userQuestionAnswers {

		var userSubQuestionAnswers []models.UserSubQuestionAnswer

		// ใช้ Preload และระบุเงื่อนไขในการ Join
		result := r.db.Model(&models.UserSubQuestionAnswer{}).
			Joins("JOIN sub_questions ON sub_questions.id = user_sub_question_answers.sub_question_id").
			Where("user_sub_question_answers.user_id = ? AND sub_questions.choice_id = ?", email, userQuestionAnswers[i].ChoiceID).
			Preload("SubQuestionChoice").
			Preload("SubQuestion").
			Limit(limit).
			Offset(offset).
			Find(&userSubQuestionAnswers)

		if result.Error != nil {
			return &[]userQuestionAnswerDto.GetUserQuestionAnswerRes{}, 0, result.Error
		}

		var subQuestionAnswersRes []userQuestionAnswerDto.SubQuestionAnswerRes
		for _, subQuestionAnswer := range userSubQuestionAnswers {
			subQuestionAnswersRes = append(subQuestionAnswersRes, userQuestionAnswerDto.SubQuestionAnswerRes{
				SubQuestionID:                        subQuestionAnswer.SubQuestionID,
				SubQuestionDescription:               subQuestionAnswer.SubQuestion.Description,
				QuestionType:                         subQuestionAnswer.SubQuestion.QuestionType,
				AnswerText:                           subQuestionAnswer.AnswerText,
				SelectedSubQuestionChoiceID:          subQuestionAnswer.SubQuestionChoice.ID,
				SelectedSubQuestionChoiceDescription: subQuestionAnswer.SubQuestionChoice.Description,
			})
		}

		userQuestionAnswerRes = append(userQuestionAnswerRes, userQuestionAnswerDto.GetUserQuestionAnswerRes{
			ID:         userQuestionAnswers[i].ID,
			QuestionID: userQuestionAnswers[i].QuestionID,
			Question: userQuestionAnswerDto.GetUserQuestionAnswersQuestion{
				Description: userQuestionAnswers[i].Question.Description,
			},
			QuestionType:     userQuestionAnswers[i].Question.QuestionType,
			ClassID:          userQuestionAnswers[i].ClassID,
			SelectedChoiceID: userQuestionAnswers[i].ChoiceID,
			SelectedChoice: userQuestionAnswerDto.GetUserQuestionAnswersChoice{
				ID:                 userQuestionAnswers[i].Choice.ID,
				Description:        userQuestionAnswers[i].Choice.Description,
				SubQuestionAnswers: subQuestionAnswersRes,
			},
			Email:      userQuestionAnswers[i].UserEmail,
			AnswerText: userQuestionAnswers[i].AnswerText,

			CreatedAt: userQuestionAnswers[i].CreatedAt,
			UpdatedAt: userQuestionAnswers[i].UpdatedAt,
		})
	}

	if result.Error != nil {
		return &[]userQuestionAnswerDto.GetUserQuestionAnswerRes{}, 0, result.Error
	}

	return &userQuestionAnswerRes, total, nil
}

func (r *userQuestionAnswerGormRepository) GetUserQuestionAnswerById(userQuestionAnswerID string) (*models.UserQuestionAnswer, error) {
	var userQuestionAnswer = new(models.UserQuestionAnswer)
	result := r.db.First(&userQuestionAnswer, "id = ?", userQuestionAnswerID)

	if result.Error != nil {
		return &models.UserQuestionAnswer{}, result.Error
	}

	if result.RowsAffected == 0 {
		return &models.UserQuestionAnswer{}, errors.New("user question answer record not found")
	}

	return userQuestionAnswer, nil
}

func (r *userQuestionAnswerGormRepository) IsUserAnsweredThisQuestion(email, questionID string) (bool, error) {
	var count int64
	result := r.db.Model(&models.UserQuestionAnswer{}).
		Where("email = ? AND question_id = ?", email, questionID).
		Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}
