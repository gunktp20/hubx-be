package usecase

import (
	"fmt"
	"strings"

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
		CreateAttendances(createAttendanceReqs []attendanceDto.CreateAttendanceReq) ([]attendanceDto.CreateAttendanceRes, error)
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

func (u *attendanceUsecase) CreateAttendances(createAttendanceReqs []attendanceDto.CreateAttendanceReq) ([]attendanceDto.CreateAttendanceRes, error) {
	var results []attendanceDto.CreateAttendanceRes
	var errMessages []string

	for _, req := range createAttendanceReqs {
		userAttendanceCount, err := u.attendanceRepo.CountAttendancesByClassSessionIDAndEmail(req.ClassSessionID, req.UserEmail)
		if err != nil {
			errMessages = append(errMessages, fmt.Sprintf("Error checking attendance for %s: %v", req.UserEmail, err))
			continue
		}
		if userAttendanceCount > 0 {
			errMessages = append(errMessages, fmt.Sprintf("User %s has already attended session %s", req.UserEmail, req.ClassSessionID))
			continue
		}

		selectedClassSession, err := u.classSessionRepo.GetClassSessionById(req.ClassSessionID)
		if err != nil {
			errMessages = append(errMessages, fmt.Sprintf("Error fetching class session %s: %v", req.ClassSessionID, err))
			continue
		}
		req.ClassID = selectedClassSession.ClassID

		isUserRegistered, err := u.classRegistrationRepo.HasUserRegisteredByClassSessionID(req.UserEmail, req.ClassSessionID)
		if err != nil {
			errMessages = append(errMessages, fmt.Sprintf("Error checking registration for %s: %v", req.UserEmail, err))
			continue
		}

		if !isUserRegistered {
			errMessages = append(errMessages, fmt.Sprintf("User %s is not registered for session %s", req.UserEmail, req.ClassSessionID))
			continue
		}

		res, err := u.attendanceRepo.CreateAttendance(&req)
		if err != nil {
			errMessages = append(errMessages, fmt.Sprintf("Error creating attendance for %s: %v", req.UserEmail, err))
			continue
		}

		results = append(results, *res)
	}

	if len(errMessages) > 0 {
		return results, fmt.Errorf("Some errors occurred: %s", strings.Join(errMessages, "; "))
	}

	return results, nil
}
