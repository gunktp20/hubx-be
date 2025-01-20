package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClassSessionStatus string

const (
	Available ClassSessionStatus = "available" // สามารถลงทะเบียนได้
	Closed    ClassSessionStatus = "closed"    // ปิดการลงทะเบียน
)

type ClassSession struct {
	ID                 string             `gorm:"type:uuid;primaryKey;not null" json:"id"`
	ClassID            string             `gorm:"type:uuid;index;not null" json:"class_id"`
	Class              Class              `gorm:"foreignKey:ClassID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"class"`
	Date               time.Time          `gorm:"index;not null" json:"date"`
	MaxCapacity        int                `gorm:"not null" json:"max_capacity"`
	ClassSessionStatus ClassSessionStatus `gorm:"type:class_session_status;not null;default:'available'" json:"class_session_status"`
	StartTime          time.Time          `gorm:"index;not null" json:"start_time"`
	EndTime            time.Time          `gorm:"index;not null" json:"end_time"`
	Location           string             `gorm:"not null" json:"location"`
	CreatedAt          time.Time
	UpdatedAt          time.Time

	Attendances []Attendance `gorm:"foreignKey:ClassSessionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"attendances"`
}

func (classSession *ClassSession) BeforeCreate(tx *gorm.DB) (err error) {
	if classSession.ID == "" {
		classSession.ID = uuid.New().String()
	}
	return
}
