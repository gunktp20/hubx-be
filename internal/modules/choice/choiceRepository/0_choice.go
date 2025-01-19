package repository

import (
	"errors"

	choiceDto "github.com/gunktp20/digital-hubx-be/internal/modules/choice/choiceDto"
	"github.com/gunktp20/digital-hubx-be/pkg/models"
)

type (
	ChoiceRepositoryService interface {
		CreateChoice(createChoiceReq *choiceDto.CreateChoiceReq) (*choiceDto.CreateChoiceRes, error)
		GetChoicesByClassID(classId string, page int, limit int) (*[]models.Choice, int64, error)
		GetChoiceById(choiceID string) (*models.Choice, error)
	}
)

func (r *choiceGormRepository) CreateChoice(createChoiceReq *choiceDto.CreateChoiceReq) (*choiceDto.CreateChoiceRes, error) {

	choice := models.Choice{
		Description: createChoiceReq.Description,
		QuestionID:  createChoiceReq.QuestionID,
	}

	if err := r.db.Create(&choice).Error; err != nil {
		return &choiceDto.CreateChoiceRes{}, err
	}

	return &choiceDto.CreateChoiceRes{
		ID:          choice.ID,
		Description: choice.Description,
		QuestionID:  choice.QuestionID,
	}, nil
}

func (r *choiceGormRepository) GetChoicesByClassID(classId string, page int, limit int) (*[]models.Choice, int64, error) {
	var choices []models.Choice
	var total int64

	query := r.db.Model(&models.Choice{})

	query.Count(&total)

	result := query.
		Preload("Choices").
		Find(&choices, "class_id = ?", classId)

	if result.Error != nil {
		return &[]models.Choice{}, 0, result.Error
	}

	return &choices, total, nil

}

func (r *choiceGormRepository) GetChoiceById(choiceID string) (*models.Choice, error) {
	var choice = new(models.Choice)
	result := r.db.First(&choice, "id = ?", choiceID)

	if result.Error != nil {
		return &models.Choice{}, result.Error
	}

	if result.RowsAffected == 0 {
		return &models.Choice{}, errors.New("choice record not found")
	}

	return choice, nil
}
