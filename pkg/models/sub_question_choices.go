package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubQuestionChoice struct {
	ID            string      `gorm:"type:uuid;primaryKey;not null" json:"id"`
	Description   string      `gorm:"not null" json:"description"`
	SubQuestionID string      `gorm:"type:uuid;index;not null;" json:"sub_question_id"`
	SubQuestion   SubQuestion `gorm:"foreignKey:SubQuestionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"sub_question"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Order int `gorm:"not null" json:"order"`

	UserSubQuestionAnswers []UserSubQuestionAnswer `gorm:"foreignKey:SubQuestionChoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user_sub_question_answers"`
}

func (subQuestionChoice *SubQuestionChoice) BeforeCreate(tx *gorm.DB) (err error) {
	if subQuestionChoice.ID == "" {
		subQuestionChoice.ID = uuid.New().String()
	}

	var maxOrder int
	if err := tx.Model(&SubQuestionChoice{}).
		Where("sub_question_id = ?", subQuestionChoice.SubQuestionID).
		Select("COALESCE(MAX(\"order\"), 0)").
		Scan(&maxOrder).Error; err != nil {
		return err
	}

	subQuestionChoice.Order = maxOrder + 1

	return
}
