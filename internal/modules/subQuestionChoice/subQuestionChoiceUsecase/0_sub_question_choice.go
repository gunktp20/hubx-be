package usecase

import (
	subQuestionRepository "github.com/gunktp20/digital-hubx-be/internal/modules/subQuestion/subQuestionRepository"
	subQuestionChoiceDto "github.com/gunktp20/digital-hubx-be/internal/modules/subQuestionChoice/subQuestionChoiceDto"
	subQuestionChoiceRepository "github.com/gunktp20/digital-hubx-be/internal/modules/subQuestionChoice/subQuestionChoiceRepository"
)

type (
	SubQuestionChoiceUsecaseService interface {
		CreateSubQuestionChoice(createSubQuestionChoiceReq *subQuestionChoiceDto.CreateSubQuestionChoicesReq) (*subQuestionChoiceDto.CreateSubQuestionChoicesRes, error)
	}

	subQuestionChoiceUsecase struct {
		subQuestionChoiceRepo subQuestionChoiceRepository.SubQuestionChoiceRepositoryService
		subQuestionRepo       subQuestionRepository.SubQuestionRepositoryService
	}
)

func NewSubQuestionChoiceUsecase(subQuestionChoiceRepo subQuestionChoiceRepository.SubQuestionChoiceRepositoryService, subQuestionRepo subQuestionRepository.SubQuestionRepositoryService) SubQuestionChoiceUsecaseService {
	return &subQuestionChoiceUsecase{subQuestionChoiceRepo: subQuestionChoiceRepo, subQuestionRepo: subQuestionRepo}
}

func (u *subQuestionChoiceUsecase) CreateSubQuestionChoice(createSubQuestionChoiceReq *subQuestionChoiceDto.CreateSubQuestionChoicesReq) (*subQuestionChoiceDto.CreateSubQuestionChoicesRes, error) {

	_, err := u.subQuestionRepo.GetSubQuestionById(createSubQuestionChoiceReq.SubQuestionID)
	if err != nil {
		return &subQuestionChoiceDto.CreateSubQuestionChoicesRes{}, err
	}

	return u.subQuestionChoiceRepo.CreateSubQuestionChoice(createSubQuestionChoiceReq)
}
