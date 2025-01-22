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
		CreateAttendances(createAttendanceReqs []attendanceDto.CreateAttendanceReq) ([]attendanceDto.CreateAttendanceRes, error)
	}
)

func (r *attendanceGormRepository) CreateAttendance(createAttendanceReq *attendanceDto.CreateAttendanceReq) (*attendanceDto.CreateAttendanceRes, error) {

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

	result := r.db.Model(&models.Attendance{}).
		Where("class_session_id = ? AND user_email = ?", classSessionID, userEmail).
		Count(&total)

	if result.Error != nil {
		return 0, result.Error
	}

	return total, nil
}

func (r *attendanceGormRepository) CreateAttendances(createAttendanceReqs []attendanceDto.CreateAttendanceReq) ([]attendanceDto.CreateAttendanceRes, error) {
	var attendances []models.Attendance
	var results []attendanceDto.CreateAttendanceRes

	for _, req := range createAttendanceReqs {
		attendances = append(attendances, models.Attendance{
			UserEmail:      req.UserEmail,
			ClassID:        req.ClassID,
			ClassSessionID: req.ClassSessionID,
		})
	}

	if err := r.db.Create(&attendances).Error; err != nil {
		return nil, err
	}

	for _, attendance := range attendances {
		results = append(results, attendanceDto.CreateAttendanceRes{
			ID:             attendance.ID,
			UserEmail:      attendance.UserEmail,
			ClassID:        attendance.ClassID,
			ClassSessionID: attendance.ClassSessionID,
		})
	}

	return results, nil
}
