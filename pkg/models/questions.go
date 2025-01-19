package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionTypes string

const (
	QChoice QuestionTypes = "choice"
	QText   QuestionTypes = "text"
)

type Question struct {
	ID           string        `gorm:"type:uuid;primaryKey;not null" json:"id"`
	Description  string        `gorm:"not null" json:"description"`
	ClassID      string        `gorm:"type:uuid;index;not null" json:"class_id"`
	Class        Class         `gorm:"foreignKey:ClassID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"class"`
	QuestionType QuestionTypes `gorm:"type:question_types;not null" json:"question_type"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Order int `gorm:"not null" json:"order"`

	Choices             []Choice             `json:"choices"`
	UserQuestionAnswers []UserQuestionAnswer `gorm:"foreignKey:QuestionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user_question_answers"`
}

func (question *Question) BeforeCreate(tx *gorm.DB) (err error) {
	if question.ID == "" {
		question.ID = uuid.New().String()
	}

	var maxOrder int
	if err := tx.Model(&Question{}).
		Where("class_id = ?", question.ClassID).
		Select("COALESCE(MAX(\"order\"), 0)").
		Scan(&maxOrder).Error; err != nil {
		return err
	}

	question.Order = maxOrder + 1

	return
}
