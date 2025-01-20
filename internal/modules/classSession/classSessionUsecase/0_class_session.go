package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/gunktp20/digital-hubx-be/external/gcs"
	classRepository "github.com/gunktp20/digital-hubx-be/internal/modules/class/classRepository"
	classRegistrationRepository "github.com/gunktp20/digital-hubx-be/internal/modules/classRegistration/classRegistrationRepository"
	classSessionDto "github.com/gunktp20/digital-hubx-be/internal/modules/classSession/classSessionDto"
	classSessionRepository "github.com/gunktp20/digital-hubx-be/internal/modules/classSession/classSessionRepository"
	"github.com/gunktp20/digital-hubx-be/pkg/config"
	"github.com/gunktp20/digital-hubx-be/pkg/utils"
)

type (
	ClassSessionUsecaseService interface {
		CreateClassSession(createClassSessionReq *classSessionDto.CreateClassSessionReq) (*classSessionDto.CreateClassSessionRes, error)
		GetAllClassSessions(class_id, class_tier string, page int, limit int) (*[]classSessionDto.ClassSessionsRes, int64, error)
		SetMaxCapacity(classSessionID string, newCapacity int) error
	}

	classSessionUsecase struct {
		cfg                   *config.Config
		classSessionRepo      classSessionRepository.ClassSessionRepositoryService
		classRepo             classRepository.ClassRepositoryService
		classRegistrationRepo classRegistrationRepository.ClassRegistrationRepositoryService
		gcsClient             gcs.GcsClientService
	}
)

func NewClassSessionUsecase(cfg *config.Config, classSessionRepo classSessionRepository.ClassSessionRepositoryService, classRepo classRepository.ClassRepositoryService, classRegistrationRepo classRegistrationRepository.ClassRegistrationRepositoryService, gcsClient gcs.GcsClientService) ClassSessionUsecaseService {
	return &classSessionUsecase{cfg: cfg, classSessionRepo: classSessionRepo, classRepo: classRepo, classRegistrationRepo: classRegistrationRepo, gcsClient: gcsClient}
}

func (u *classSessionUsecase) CreateClassSession(createClassSessionReq *classSessionDto.CreateClassSessionReq) (*classSessionDto.CreateClassSessionRes, error) {

	selectedClass, err := u.classRepo.GetClassById(createClassSessionReq.ClassID)
	if err != nil {
		return &classSessionDto.CreateClassSessionRes{}, err
	}

	dateFields := []time.Time{
		createClassSessionReq.Date, createClassSessionReq.StartTime, createClassSessionReq.EndTime,
	}

	// ? Check all fields about date and time is future
	_, err = utils.AreAllFutureDates(dateFields)
	if err != nil {
		return &classSessionDto.CreateClassSessionRes{}, err
	}

	// ? Check is date conflict ?
	isDateConflict, err := u.classSessionRepo.CheckSessionDateConflict(createClassSessionReq.ClassID, string(selectedClass.ClassTier), createClassSessionReq.Date)

	// ? Calculate cancellation deadline
	var cancellationDeadline = createClassSessionReq.Date.AddDate(0, 0, -7)

	if err != nil {
		return &classSessionDto.CreateClassSessionRes{}, err
	}

	if isDateConflict {
		return &classSessionDto.CreateClassSessionRes{}, errors.New("class session date conflicts with an existing session")
	}

	if createClassSessionReq.MaxCapacity <= 0 {
		return &classSessionDto.CreateClassSessionRes{}, errors.New("max capacity must be greater than zero")
	}

	if createClassSessionReq.MaxCapacity > u.cfg.BusinessLogic.MaxCapacityPerSession {
		return &classSessionDto.CreateClassSessionRes{}, fmt.Errorf("capacity exceeds the maximum allowed limit of %d for a session", u.cfg.BusinessLogic.MaxCapacityPerSession)
	}

	return u.classSessionRepo.CreateClassSession(createClassSessionReq, cancellationDeadline)
}

func (u *classSessionUsecase) GetAllClassSessions(class_id, class_tier string, page int, limit int) (*[]classSessionDto.ClassSessionsRes, int64, error) {
	classSessiones, total, err := u.classSessionRepo.GetAllClassSessions(class_id, class_tier, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return classSessiones, total, nil
}

func (u *classSessionUsecase) SetMaxCapacity(classSessionID string, newCapacity int) error {

	countRegistrations, err := u.classRegistrationRepo.CountRegistrationWithClassSessionID(classSessionID)
	if err != nil {
		return err
	}

	if newCapacity > u.cfg.BusinessLogic.MaxCapacityPerSession {
		return fmt.Errorf("capacity exceeds the maximum allowed limit of %d for a session", u.cfg.BusinessLogic.MaxCapacityPerSession)
	}

	if countRegistrations > newCapacity {
		return fmt.Errorf("new capacity (%d) is less than the current number of registrations (%d)", newCapacity, countRegistrations)
	}

	err = u.classSessionRepo.SetMaxCapacity(classSessionID, newCapacity)
	if err != nil {
		return err
	}

	return nil
}
