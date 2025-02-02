package usecase

import (
	"errors"
	"fmt"
	"strings"
	"time"

	classRepository "github.com/gunktp20/digital-hubx-be/internal/modules/class/classRepository"
	classRegistrationDto "github.com/gunktp20/digital-hubx-be/internal/modules/classRegistration/classRegistrationDto"
	classRegistrationRepository "github.com/gunktp20/digital-hubx-be/internal/modules/classRegistration/classRegistrationRepository"
	classSessionRepository "github.com/gunktp20/digital-hubx-be/internal/modules/classSession/classSessionRepository"
	questionRepository "github.com/gunktp20/digital-hubx-be/internal/modules/question/questionRepository"
	userQuestionAnswerRepository "github.com/gunktp20/digital-hubx-be/internal/modules/userQuestionAnswer/userQuestionAnswerRepository"
	"github.com/gunktp20/digital-hubx-be/pkg/config"
	"github.com/gunktp20/digital-hubx-be/pkg/utils"
)

type (
	ClassRegistrationUsecaseService interface {
		CreateClassRegistration(createClassRegistrationReq *classRegistrationDto.CreateClassRegistrationReq, email string) (*classRegistrationDto.CreateClassRegistrationRes, error)
		GetUserRegistrations(email string, page int, limit int) (*[]classRegistrationDto.GetUserRegistrationsRes, int64, error)
		CancelClassRegistration(email, classID string) error
		ResetCancelledQuota(resetCancelledQuotaReq *classRegistrationDto.ResetCancelledQuotaReq) error
		DeleteUserClassRegistrationBySession(userEmail, classSessionID string) error
	}

	classRegistrationUsecase struct {
		classRegistrationRepo  classRegistrationRepository.ClassRegistrationRepositoryService
		classSessionRepo       classSessionRepository.ClassSessionRepositoryService
		classRepo              classRepository.ClassRepositoryService
		userQuestionAnswerRepo userQuestionAnswerRepository.UserQuestionAnswerRepositoryService
		questionRepo           questionRepository.QuestionRepositoryService
		conf                   config.Config
	}
)

func NewClassRegistrationUsecase(
	conf *config.Config,
	classRegistrationRepo classRegistrationRepository.ClassRegistrationRepositoryService,
	classSessionRepo classSessionRepository.ClassSessionRepositoryService,
	classRepo classRepository.ClassRepositoryService,
	userQuestionAnswerRepo userQuestionAnswerRepository.UserQuestionAnswerRepositoryService,
	questionRepo questionRepository.QuestionRepositoryService,
) ClassRegistrationUsecaseService {
	return &classRegistrationUsecase{
		conf:                   *conf,
		classRegistrationRepo:  classRegistrationRepo,
		classSessionRepo:       classSessionRepo,
		classRepo:              classRepo,
		userQuestionAnswerRepo: userQuestionAnswerRepo,
		questionRepo:           questionRepo,
	}
}

func (u *classRegistrationUsecase) CreateClassRegistration(createClassRegistrationReq *classRegistrationDto.CreateClassRegistrationReq, email string) (*classRegistrationDto.CreateClassRegistrationRes, error) {

	classSession, err := u.classSessionRepo.GetClassSessionById(createClassRegistrationReq.ClassSessionID)
	if err != nil {
		return &classRegistrationDto.CreateClassRegistrationRes{}, err
	}

	if classSession.ClassID != createClassRegistrationReq.ClassID {
		return &classRegistrationDto.CreateClassRegistrationRes{}, errors.New("actual class id of class session doesn't match with class that you provided")
	}

	// ? Verify that the class has an is_remove status of false
	selectedClass, err := u.classRepo.GetClassById(createClassRegistrationReq.ClassID)
	if err != nil {
		return &classRegistrationDto.CreateClassRegistrationRes{}, err
	}

	if selectedClass.IsRemove {
		return &classRegistrationDto.CreateClassRegistrationRes{}, errors.New("registration is not allowed for this class as it has been removed")
	}

	//  ? Check if the user is already registered for the class session.
	isRegistered, err := u.classRegistrationRepo.HasUserRegistered(email, createClassRegistrationReq.ClassID)
	if err != nil {
		return &classRegistrationDto.CreateClassRegistrationRes{}, err
	}
	if isRegistered {
		return &classRegistrationDto.CreateClassRegistrationRes{}, errors.New("user has already registered for this class session")
	}

	// ? Check if the registration has reached the maximum capacity
	maxCapacity, err := u.classSessionRepo.GetMaxCapacityOfClassSessionById(createClassRegistrationReq.ClassSessionID)
	if err != nil {
		return &classRegistrationDto.CreateClassRegistrationRes{}, err
	}
	totalRegistrations, err := u.classRegistrationRepo.CountRegistrationWithClassSessionID(createClassRegistrationReq.ClassSessionID)
	if err != nil {
		return &classRegistrationDto.CreateClassRegistrationRes{}, err
	}
	if totalRegistrations >= maxCapacity {
		return &classRegistrationDto.CreateClassRegistrationRes{}, errors.New("registration has reached the maximum capacity")
	}

	// ? Check is event date valid for register
	eventDateValidForReg, err := utils.IsEventDateValidForReg(classSession.Date)
	if err != nil {
		return &classRegistrationDto.CreateClassRegistrationRes{}, err
	}
	if !eventDateValidForReg {
		return &classRegistrationDto.CreateClassRegistrationRes{}, errors.New("registration not allowed for this date")
	}

	cancelledCount, err := u.classRegistrationRepo.CountUserCancelledRegistrationsByEmail(email, createClassRegistrationReq.ClassID)
	if err != nil {
		return &classRegistrationDto.CreateClassRegistrationRes{}, err
	}

	if cancelledCount >= u.conf.BusinessLogic.MaxCancelPerClass {
		return &classRegistrationDto.CreateClassRegistrationRes{}, errors.New("you cannot register for this class because you have reached the maximum cancellation limit. Please contact the administrator if you believe this is an error")
	}

	if selectedClass.EnableQuestion {
		questionsCount, err := u.questionRepo.CountQuestionsByClassID(createClassRegistrationReq.ClassID)
		if err != nil {
			return &classRegistrationDto.CreateClassRegistrationRes{}, err
		}

		if questionsCount > 0 {
			answersCount, err := u.userQuestionAnswerRepo.CountUserAnswersByEmailAndClassId(email, createClassRegistrationReq.ClassID)
			if err != nil {
				return &classRegistrationDto.CreateClassRegistrationRes{}, err
			}

			if answersCount <= 0 {
				return &classRegistrationDto.CreateClassRegistrationRes{}, fmt.Errorf("please complete the survey before registering for this class")
			}
		}
	}

	return u.classRegistrationRepo.CreateClassRegistration(createClassRegistrationReq, email)
}

func (u *classRegistrationUsecase) GetUserRegistrations(email string, page int, limit int) (*[]classRegistrationDto.GetUserRegistrationsRes, int64, error) {

	userClassRegistration, total, err := u.classRegistrationRepo.GetUserRegistrations(email, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return userClassRegistration, total, nil
}

func (u *classRegistrationUsecase) CancelClassRegistration(email, classID string) error {
	// ? Fetch the user's registration for the given classID, including the associated ClassSession
	registration, err := u.classRegistrationRepo.GetUserClassRegistrationByClassIDAndEmail(classID, email)
	if err != nil {
		// ? Return an error if no registration is found or another issue occurs
		return fmt.Errorf("failed to fetch registration details: %v", err)
	}

	// ? Calculate the number of days remaining until the class starts
	daysUntilClass := int(time.Until(registration.ClassSession.Date).Hours() / 24)

	// ? Check if cancellation is allowed based on business logic
	if daysUntilClass < u.conf.BusinessLogic.DaysBeforeClassStartForCancellation {
		return fmt.Errorf("cancellation is not allowed within %d days before the class start date", u.conf.BusinessLogic.DaysBeforeClassStartForCancellation)
	}

	// ? Perform the cancellation by updating the registration status
	err = u.classRegistrationRepo.CancelClassRegistration(email, classID)
	if err != nil {
		// ? Return an error if the cancellation fails
		return fmt.Errorf("failed to cancel class registration: %v", err)
	}

	return nil
}

func (u *classRegistrationUsecase) ResetCancelledQuota(resetCancelledQuotaReq *classRegistrationDto.ResetCancelledQuotaReq) error {

	resetCancelledQuotaReq.UserEmail = strings.ToLower(resetCancelledQuotaReq.UserEmail)

	err := u.classRegistrationRepo.ResetCancelledQuota(resetCancelledQuotaReq)
	if err != nil {
		return err
	}

	return nil
}

func (u *classRegistrationUsecase) DeleteUserClassRegistrationBySession(userEmail, classSessionID string) error {
	err := u.classRegistrationRepo.DeleteUserClassRegistrationBySession(userEmail, classSessionID)
	if err != nil {
		return fmt.Errorf("failed to delete user class registration: %v", err)
	}

	return nil
}
