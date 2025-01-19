package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubQuestion struct {
	ID          string   `gorm:"type:uuid;primaryKey;not null" json:"id"`
	Description string   `gorm:"not null" json:"description"`                 // เพิ่ม
	QuestionID  string   `gorm:"type:uuid;index;not null" json:"question_id"` // เพิ่ม
	Question    Question `gorm:"foreignKey:QuestionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"question"`
	ChoiceID    string   `gorm:"type:uuid;index;not null" json:"choice_id"` // เพิ่ม
	Choice      Choice   `gorm:"foreignKey:ChoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"choice"`

	QuestionType QuestionTypes `gorm:"type:question_types;not null" json:"question_type"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Order int `gorm:"not null" json:"order"`

	SubQuestionChoices     []SubQuestionChoice     `gorm:"foreignKey:SubQuestionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"sub_question_choices"`
	UserSubQuestionAnswers []UserSubQuestionAnswer `gorm:"foreignKey:SubQuestionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user_sub_question_answers"`
}

func (subQuestion *SubQuestion) BeforeCreate(tx *gorm.DB) (err error) {
	if subQuestion.ID == "" {
		subQuestion.ID = uuid.New().String()
	}

	var maxOrder int
	if err := tx.Model(&SubQuestion{}).
		Where("choice_id = ?", subQuestion.ChoiceID).
		Select("COALESCE(MAX(\"order\"), 0)").
		Scan(&maxOrder).Error; err != nil {
		return err
	}

	subQuestion.Order = maxOrder + 1

	return
}
