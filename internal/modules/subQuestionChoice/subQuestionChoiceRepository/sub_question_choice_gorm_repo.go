package repository

import (
	"gorm.io/gorm"
)

type subQuestionChoiceGormRepository struct {
	db *gorm.DB
}

func NewSubQuestionChoiceGormRepository(db *gorm.DB) SubQuestionChoiceRepositoryService {
	return &subQuestionChoiceGormRepository{db}
}
