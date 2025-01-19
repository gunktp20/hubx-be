package repository

import (
	"errors"

	subQuestionDto "github.com/gunktp20/digital-hubx-be/internal/modules/subQuestion/subQuestionDto"
	"github.com/gunktp20/digital-hubx-be/pkg/models"
)

type (
	SubQuestionRepositoryService interface {
		CreateSubQuestion(createSubQuestionReq *subQuestionDto.CreateSubQuestionReq) (*subQuestionDto.CreateSubQuestionRes, error)
		GetSubQuestionsByQuestionID(questionID string, page int, limit int) (*[]models.SubQuestion, int64, error)
		GetSubQuestionsByChoiceID(choiceID string, page int, limit int) (*[]models.SubQuestion, int64, error)
		GetSubQuestionById(subQuestionID string) (*models.SubQuestion, error)
	}
)

func (r *subQuestionGormRepository) CreateSubQuestion(createSubQuestionReq *subQuestionDto.CreateSubQuestionReq) (*subQuestionDto.CreateSubQuestionRes, error) {

	subQuestion := models.SubQuestion{
		Description:  createSubQuestionReq.Description,
		ChoiceID:     createSubQuestionReq.ChoiceID,
		QuestionID:   createSubQuestionReq.QuestionID,
		QuestionType: models.QuestionTypes(createSubQuestionReq.QuestionType),
	}

	if err := r.db.Create(&subQuestion).Error; err != nil {
		return &subQuestionDto.CreateSubQuestionRes{}, err
	}

	if err := r.db.Preload("Choice").Preload("Question").First(&subQuestion).Error; err != nil {
		return &subQuestionDto.CreateSubQuestionRes{}, err
	}

	return &subQuestionDto.CreateSubQuestionRes{
		ID:          subQuestion.ID,
		Description: subQuestion.Description,
		QuestionID:  subQuestion.QuestionID,
		ChoiceID:    subQuestion.ChoiceID,
		Choice: subQuestionDto.ChoiceRes{
			Description: subQuestion.Choice.Description,
		},

		QuestionType: subQuestion.QuestionType,
		CreatedAt:    subQuestion.CreatedAt,
		UpdatedAt:    subQuestion.UpdatedAt,
	}, nil

}

func (r *subQuestionGormRepository) GetSubQuestionsByQuestionID(questionID string, page int, limit int) (*[]models.SubQuestion, int64, error) {
	var subQuestions []models.SubQuestion
	var total int64

	query := r.db.Model(&models.SubQuestion{})

	query.Count(&total)

	result := query.
		Preload("Choice").
		Find(&subQuestions, "question_id = ?", questionID)

	if result.Error != nil {
		return &[]models.SubQuestion{}, 0, result.Error
	}

	return &subQuestions, total, nil

}

func (r *subQuestionGormRepository) GetSubQuestionsByChoiceID(choiceID string, page int, limit int) (*[]models.SubQuestion, int64, error) {
	var subQuestions []models.SubQuestion
	var total int64

	query := r.db.Model(&models.SubQuestion{})

	query.Count(&total)

	result := query.
		Preload("Choice").
		Find(&subQuestions, "choice_id = ?", choiceID)

	if result.Error != nil {
		return &[]models.SubQuestion{}, 0, result.Error
	}

	return &subQuestions, total, nil

}

func (r *subQuestionGormRepository) GetSubQuestionById(subQuestionID string) (*models.SubQuestion, error) {
	var subQuestion = new(models.SubQuestion)
	result := r.db.First(&subQuestion, "id = ?", subQuestionID)

	if result.Error != nil {
		return &models.SubQuestion{}, result.Error
	}

	if result.RowsAffected == 0 {
		return &models.SubQuestion{}, errors.New("subQuestion record not found")
	}

	return subQuestion, nil
}
