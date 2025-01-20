package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Attendance struct {
	ID             string       `gorm:"type:uuid;primaryKey;not null" json:"id"`
	UserEmail      string       `gorm:"type:varchar(255);index;not null;" json:"user_email"`
	ClassID        string       `gorm:"type:uuid;index;not null;" json:"class_id"`
	Class          Class        `gorm:"foreignKey:ClassID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"class"`
	ClassSessionID string       `gorm:"type:uuid;index;not null;" json:"class_session_id"`
	ClassSession   ClassSession `gorm:"foreignKey:ClassSessionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"class_session"`
	CheckInTime    time.Time    `gorm:"autoCreateTime" json:"check_in_time"`
	Status         string       `gorm:"type:varchar(50);default:'Present';not null" json:"status"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (attendance *Attendance) BeforeCreate(tx *gorm.DB) (err error) {
	if attendance.ID == "" {
		attendance.ID = uuid.New().String()
	}
	return
}
