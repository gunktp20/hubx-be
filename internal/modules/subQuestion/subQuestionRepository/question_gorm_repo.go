package repository

import (
	"gorm.io/gorm"
)

type subQuestionGormRepository struct {
	db *gorm.DB
}

func NewSubQuestionGormRepository(db *gorm.DB) SubQuestionRepositoryService {
	return &subQuestionGormRepository{db}
}
