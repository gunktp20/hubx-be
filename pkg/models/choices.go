package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Choice struct {
	ID          string   `gorm:"type:uuid;primaryKey;not null" json:"id"`
	Description string   `gorm:"not null" json:"description"`
	QuestionID  string   `gorm:"type:uuid;index;not null;" json:"question_id"` // ?
	Question    Question `gorm:"foreignKey:QuestionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"question"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Order int `gorm:"not null" json:"order"`

	SubQuestions        []SubQuestion        `gorm:"foreignKey:ChoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"sub_questions"`
	UserQuestionAnswers []UserQuestionAnswer `gorm:"foreignKey:ChoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user_question_answers"`
}

func (choice *Choice) BeforeCreate(tx *gorm.DB) (err error) {
	if choice.ID == "" {
		choice.ID = uuid.New().String()
	}

	var maxOrder int
	if err := tx.Model(&Choice{}).
		Where("question_id = ?", choice.QuestionID).
		Select("COALESCE(MAX(\"order\"), 0)").
		Scan(&maxOrder).Error; err != nil {
		return err
	}

	choice.Order = maxOrder + 1

	return
}
