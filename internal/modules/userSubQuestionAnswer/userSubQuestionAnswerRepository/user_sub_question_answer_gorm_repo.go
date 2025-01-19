package repository

import (
	"gorm.io/gorm"
)

type userSubQuestionAnswerGormRepository struct {
	db *gorm.DB
}

func NewUserSubQuestionAnswerGormRepository(db *gorm.DB) UserSubQuestionAnswerRepositoryService {
	return &userSubQuestionAnswerGormRepository{db}
}
