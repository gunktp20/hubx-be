package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSubQuestionAnswer struct {
	ID                  string  `gorm:"type:uuid;primaryKey;uniqueIndex;not null" json:"id"`
	UserEmail           string  `gorm:"type:varchar(255);not null;" json:"user_email"`
	SubQuestionID       string  `gorm:"type:uuid;not null;" json:"sub_question_id"`
	SubQuestionChoiceID *string `gorm:"type:uuid;null;" json:"sub_question_choice_id"`
	ClassID             string  `gorm:"type:uuid;not null;" json:"class_id"`
	AnswerText          string  `gorm:"null" json:"answer_text"`

	User              User              `gorm:"foreignKey:UserEmail;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Class             Class             `gorm:"foreignKey:ClassID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"class"`
	SubQuestion       SubQuestion       `gorm:"foreignKey:SubQuestionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"sub_question"`
	SubQuestionChoice SubQuestionChoice `gorm:"foreignKey:SubQuestionChoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"sub_question_choice"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (userSubQuestionAnswer *UserSubQuestionAnswer) BeforeCreate(tx *gorm.DB) (err error) {
	if userSubQuestionAnswer.ID == "" {
		userSubQuestionAnswer.ID = uuid.New().String()
	}
	return
}
