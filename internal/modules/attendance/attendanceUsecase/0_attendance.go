package usecase

import (
	"fmt"

	attendanceDto "github.com/gunktp20/digital-hubx-be/internal/modules/attendance/attendanceDto"
	attendanceRepository "github.com/gunktp20/digital-hubx-be/internal/modules/attendance/attendanceRepository"
	classRegistrationRepository "github.com/gunktp20/digital-hubx-be/internal/modules/classRegistration/classRegistrationRepository"
	classSessionRepository "github.com/gunktp20/digital-hubx-be/internal/modules/classSession/classSessionRepository"
	"github.com/gunktp20/digital-hubx-be/pkg/models"
)

type (
	AttendanceUsecaseService interface {
		CreateAttendance(createAttendanceReq *attendanceDto.CreateAttendanceReq) (*attendanceDto.CreateAttendanceRes, error)
		GetAttendancesByClassID(classId string, page int, limit int) (*[]models.Attendance, int64, error)
	}

	attendanceUsecase struct {
		attendanceRepo        attendanceRepository.AttendanceRepositoryService
		classSessionRepo      classSessionRepository.ClassSessionRepositoryService
		classRegistrationRepo classRegistrationRepository.ClassRegistrationRepositoryService
	}
)

func NewAttendanceUsecase(attendanceRepo attendanceRepository.AttendanceRepositoryService, classSessionRepo classSessionRepository.ClassSessionRepositoryService, classRegistrationRepo classRegistrationRepository.ClassRegistrationRepositoryService) AttendanceUsecaseService {
	return &attendanceUsecase{attendanceRepo: attendanceRepo, classSessionRepo: classSessionRepo, classRegistrationRepo: classRegistrationRepo}
}

func (u *attendanceUsecase) CreateAttendance(createAttendanceReq *attendanceDto.CreateAttendanceReq) (*attendanceDto.CreateAttendanceRes, error) {

	userAttendanceCount, err := u.attendanceRepo.CountAttendancesByClassSessionIDAndEmail(createAttendanceReq.ClassSessionID, createAttendanceReq.UserEmail)
	if err != nil {
		return &attendanceDto.CreateAttendanceRes{}, err
	}
	if userAttendanceCount > 0 {
		return &attendanceDto.CreateAttendanceRes{}, fmt.Errorf("user with email %s has already attended class session with ID %s", createAttendanceReq.UserEmail, createAttendanceReq.ClassSessionID)
	}

	selectedClassSession, err := u.classSessionRepo.GetClassSessionById(createAttendanceReq.ClassSessionID)
	if err != nil {
		return &attendanceDto.CreateAttendanceRes{}, err
	}
	createAttendanceReq.ClassID = selectedClassSession.ClassID

	isUserRegistered, err := u.classRegistrationRepo.HasUserRegisteredByClassSessionID(createAttendanceReq.UserEmail, createAttendanceReq.ClassSessionID)
	if err != nil {
		return &attendanceDto.CreateAttendanceRes{}, err
	}

	if !isUserRegistered {
		return &attendanceDto.CreateAttendanceRes{}, fmt.Errorf("user with email %s has not registered for class session with ID %s", createAttendanceReq.UserEmail, createAttendanceReq.ClassSessionID)
	}

	return u.attendanceRepo.CreateAttendance(createAttendanceReq)
}

func (u *attendanceUsecase) GetAttendancesByClassID(classId string, page int, limit int) (*[]models.Attendance, int64, error) {

	res, total, err := u.attendanceRepo.GetAttendancesByClassID(classId, page, limit)
	if err != nil {
		return &[]models.Attendance{}, total, nil
	}

	return res, total, nil
}
