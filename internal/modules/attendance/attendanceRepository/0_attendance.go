package repository

import (
	"errors"

	attendanceDto "github.com/gunktp20/digital-hubx-be/internal/modules/attendance/attendanceDto"
	"github.com/gunktp20/digital-hubx-be/pkg/models"
)

type (
	AttendanceRepositoryService interface {
		CreateAttendance(createAttendanceReq *attendanceDto.CreateAttendanceReq) (*attendanceDto.CreateAttendanceRes, error)
		GetAttendancesByClassID(classID string, page int, limit int) (*[]models.Attendance, int64, error)
		GetAttendancesByClassSessionID(classSessionID string, page int, limit int) (*[]models.Attendance, int64, error)
		GetAttendanceById(attendanceID string) (*models.Attendance, error)
		CountAttendancesByClassSessionIDAndEmail(classSessionID string, userEmail string) (int64, error)
	}
)

func (r *attendanceGormRepository) CreateAttendance(createAttendanceReq *attendanceDto.CreateAttendanceReq) (*attendanceDto.CreateAttendanceRes, error) {

	// ID             string    `gorm:"type:uuid;primaryKey;not null" json:"id"`
	// UserEmail      string    `gorm:"type:varchar(255);index;not null;" json:"user_email"`
	// ClassID        string    `gorm:"type:uuid;index;not null;" json:"class_id"`
	// Class          Class     `gorm:"foreignKey:ClassID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"class"`
	// ClassSessionID string    `gorm:"type:uuid;index;not null;" json:"class_session_id"`
	// ClassSession   Question  `gorm:"foreignKey:ClassSessionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"class_session"`
	// CheckInTime    time.Time `gorm:"autoCreateTime" json:"check_in_time"`
	// Status         string    `gorm:"type:varchar(50);default:'Present';not null" json:"status"`

	attendance := models.Attendance{
		UserEmail:      createAttendanceReq.UserEmail,
		ClassID:        createAttendanceReq.ClassID,
		ClassSessionID: createAttendanceReq.ClassSessionID,
	}

	if err := r.db.Create(&attendance).Error; err != nil {
		return &attendanceDto.CreateAttendanceRes{}, err
	}

	return &attendanceDto.CreateAttendanceRes{
		ID:             attendance.ID,
		UserEmail:      attendance.UserEmail,
		ClassID:        attendance.ClassID,
		ClassSessionID: attendance.ClassSessionID,
	}, nil
}

func (r *attendanceGormRepository) GetAttendancesByClassID(classID string, page int, limit int) (*[]models.Attendance, int64, error) {
	var attendances []models.Attendance
	var total int64

	query := r.db.Model(&models.Attendance{})

	query.Count(&total)

	result := query.
		Find(&attendances, "class_id = ?", classID)

	if result.Error != nil {
		return &[]models.Attendance{}, 0, result.Error
	}

	return &attendances, total, nil

}

func (r *attendanceGormRepository) GetAttendancesByClassSessionID(classSessionID string, page int, limit int) (*[]models.Attendance, int64, error) {
	var attendances []models.Attendance
	var total int64

	query := r.db.Model(&models.Attendance{})

	query.Count(&total)

	result := query.
		Find(&attendances, "class_session_id = ?", classSessionID)

	if result.Error != nil {
		return &[]models.Attendance{}, 0, result.Error
	}

	return &attendances, total, nil

}

func (r *attendanceGormRepository) GetAttendanceById(attendanceID string) (*models.Attendance, error) {
	var attendance = new(models.Attendance)
	result := r.db.First(&attendance, "id = ?", attendanceID)

	if result.Error != nil {
		return &models.Attendance{}, result.Error
	}

	if result.RowsAffected == 0 {
		return &models.Attendance{}, errors.New("attendance record not found")
	}

	return attendance, nil
}

func (r *attendanceGormRepository) CountAttendancesByClassSessionIDAndEmail(classSessionID string, userEmail string) (int64, error) {
	var total int64

	// นับจำนวน attendances ที่ตรงกับ class_session_id และ email
	result := r.db.Model(&models.Attendance{}).
		Where("class_session_id = ? AND user_email = ?", classSessionID, userEmail).
		Count(&total)

	if result.Error != nil {
		return 0, result.Error
	}

	return total, nil
}
