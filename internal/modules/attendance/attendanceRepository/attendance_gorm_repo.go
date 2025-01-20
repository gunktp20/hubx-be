package repository

import (
	"gorm.io/gorm"
)

type attendanceGormRepository struct {
	db *gorm.DB
}

func NewAttendanceGormRepository(db *gorm.DB) AttendanceRepositoryService {
	return &attendanceGormRepository{db}
}
