package repository

import (
	"gorm.io/gorm"
)

type userQuestionAnswerGormRepository struct {
	db *gorm.DB
}

func NewUserQuestionAnswerGormRepository(db *gorm.DB) UserQuestionAnswerRepositoryService {
	return &userQuestionAnswerGormRepository{db}
}
