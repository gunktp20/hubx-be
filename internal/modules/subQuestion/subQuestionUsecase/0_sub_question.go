package usecase

import (
	choiceRepository "github.com/gunktp20/digital-hubx-be/internal/modules/choice/choiceRepository"
	classRepository "github.com/gunktp20/digital-hubx-be/internal/modules/class/classRepository"
	subQuestionDto "github.com/gunktp20/digital-hubx-be/internal/modules/subQuestion/subQuestionDto"
	subQuestionRepository "github.com/gunktp20/digital-hubx-be/internal/modules/subQuestion/subQuestionRepository"
	"github.com/gunktp20/digital-hubx-be/pkg/models"
)

type (
	SubQuestionUsecaseService interface {
		CreateSubQuestion(createSubQuestionReq *subQuestionDto.CreateSubQuestionReq) (*subQuestionDto.CreateSubQuestionRes, error)
		GetSubQuestionsByQuestionID(questionID string, page int, limit int) (*[]models.SubQuestion, int64, error)
		GetSubQuestionsByChoiceID(choiceID string, page int, limit int) (*[]models.SubQuestion, int64, error)
	}

	subQuestionUsecase struct {
		subQuestionRepo subQuestionRepository.SubQuestionRepositoryService
		choiceRepo      choiceRepository.ChoiceRepositoryService
		classRepo       classRepository.ClassRepositoryService
		// gcsClient gcs.GcsClientService
	}
)

func NewSubQuestionUsecase(subQuestionRepo subQuestionRepository.SubQuestionRepositoryService, classRepo classRepository.ClassRepositoryService, choiceRepo choiceRepository.ChoiceRepositoryService) SubQuestionUsecaseService {
	return &subQuestionUsecase{subQuestionRepo: subQuestionRepo, classRepo: classRepo, choiceRepo: choiceRepo}
}

func (u *subQuestionUsecase) CreateSubQuestion(createSubQuestionReq *subQuestionDto.CreateSubQuestionReq) (*subQuestionDto.CreateSubQuestionRes, error) {

	selectedChoice, err := u.choiceRepo.GetChoiceById(createSubQuestionReq.ChoiceID)
	if err != nil {
		return &subQuestionDto.CreateSubQuestionRes{}, err
	}

	createSubQuestionReq.QuestionID = selectedChoice.QuestionID

	return u.subQuestionRepo.CreateSubQuestion(createSubQuestionReq)
}

func (u *subQuestionUsecase) GetSubQuestionsByQuestionID(questionID string, page int, limit int) (*[]models.SubQuestion, int64, error) {

	res, total, err := u.subQuestionRepo.GetSubQuestionsByQuestionID(questionID, page, limit)
	if err != nil {
		return &[]models.SubQuestion{}, total, nil
	}

	return res, total, nil
}

func (u *subQuestionUsecase) GetSubQuestionsByChoiceID(choiceID string, page int, limit int) (*[]models.SubQuestion, int64, error) {

	res, total, err := u.subQuestionRepo.GetSubQuestionsByChoiceID(choiceID, page, limit)
	if err != nil {
		return &[]models.SubQuestion{}, total, nil
	}

	return res, total, nil
}
