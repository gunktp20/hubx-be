package models

import (
	"time"
)

type User struct {
	// ID        string `gorm:"type:uuid;primaryKey;not null" json:"id"`
	// Email     string `gorm:"primaryKey;not null" json:"email"`
	Email     string `gorm:"primaryKey;type:varchar(255);not null" json:"email"`
	CreatedAt time.Time
	UpdatedAt time.Time

	UserClassRegistrations []UserClassRegistration `gorm:"foreignKey:UserEmail;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user_class_registrations"`
	UserQuestionAnswers    []UserQuestionAnswer    `gorm:"foreignKey:UserEmail;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user_question_answers"`
	UserSubQuestionAnswers []UserSubQuestionAnswer `gorm:"foreignKey:UserEmail;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user_sub_question_answers"`
}
