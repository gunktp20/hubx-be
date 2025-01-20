package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserQuestionAnswer struct {
	ID         string  `gorm:"type:uuid;primaryKey;not null" json:"id"`
	UserEmail  string  `gorm:"type:varchar(255);index;not null" json:"user_email"`
	QuestionID string  `gorm:"type:uuid;index;not null;" json:"question_id"`
	ClassID    string  `gorm:"type:uuid;index;not null;" json:"class_id"`
	ChoiceID   *string `gorm:"type:uuid;null;" json:"choice_id"`
	AnswerText *string `gorm:"null" json:"answer_text"`

	// User     User     `gorm:"foreignKey:UserEmail;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Class    Class    `gorm:"foreignKey:ClassID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"class"`
	Question Question `gorm:"foreignKey:QuestionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"question"`
	Choice   Choice   `gorm:"foreignKey:ChoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"choice"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (userQuestionAnswer *UserQuestionAnswer) BeforeCreate(tx *gorm.DB) (err error) {
	if userQuestionAnswer.ID == "" {
		userQuestionAnswer.ID = uuid.New().String()
	}
	return
}
