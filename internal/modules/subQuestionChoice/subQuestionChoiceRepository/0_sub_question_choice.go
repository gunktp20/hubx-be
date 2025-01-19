package repository

import (
	subQuestionChoiceDto "github.com/gunktp20/digital-hubx-be/internal/modules/subQuestionChoice/subQuestionChoiceDto"
	"github.com/gunktp20/digital-hubx-be/pkg/models"
)

type (
	SubQuestionChoiceRepositoryService interface {
		CreateSubQuestionChoice(createSubQuestionChoiceReq *subQuestionChoiceDto.CreateSubQuestionChoicesReq) (*subQuestionChoiceDto.CreateSubQuestionChoicesRes, error)
	}
)

func (r *subQuestionChoiceGormRepository) CreateSubQuestionChoice(createSubQuestionChoiceReq *subQuestionChoiceDto.CreateSubQuestionChoicesReq) (*subQuestionChoiceDto.CreateSubQuestionChoicesRes, error) {

	subQuestionChoice := models.SubQuestionChoice{
		Description:   createSubQuestionChoiceReq.Description,
		SubQuestionID: createSubQuestionChoiceReq.SubQuestionID,
	}

	if err := r.db.Create(&subQuestionChoice).Error; err != nil {
		return &subQuestionChoiceDto.CreateSubQuestionChoicesRes{}, err
	}

	return &subQuestionChoiceDto.CreateSubQuestionChoicesRes{
		ID:            subQuestionChoice.ID,
		Description:   subQuestionChoice.Description,
		SubQuestionID: subQuestionChoice.SubQuestionID,
		CreatedAt:     subQuestionChoice.CreatedAt,
		UpdatedAt:     subQuestionChoice.UpdatedAt,
	}, nil
}
