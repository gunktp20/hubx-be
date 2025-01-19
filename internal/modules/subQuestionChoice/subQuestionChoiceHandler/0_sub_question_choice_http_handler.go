package handler

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	subQuestionChoiceDto "github.com/gunktp20/digital-hubx-be/internal/modules/subQuestionChoice/subQuestionChoiceDto"
	subQuestionChoiceUsecase "github.com/gunktp20/digital-hubx-be/internal/modules/subQuestionChoice/subQuestionChoiceUsecase"
	"github.com/gunktp20/digital-hubx-be/pkg/response"
	"github.com/gunktp20/digital-hubx-be/pkg/utils"
)

type (
	SubQuestionChoiceHttpHandlerService interface {
		CreateSubQuestionChoice(c *fiber.Ctx) error
	}

	subQuestionChoiceHttpHandler struct {
		subQuestionChoiceUsecase subQuestionChoiceUsecase.SubQuestionChoiceUsecaseService
	}
)

func NewSubQuestionChoiceHttpHandler(usecase subQuestionChoiceUsecase.SubQuestionChoiceUsecaseService) SubQuestionChoiceHttpHandlerService {
	return &subQuestionChoiceHttpHandler{subQuestionChoiceUsecase: usecase}
}

func (h *subQuestionChoiceHttpHandler) CreateSubQuestionChoice(c *fiber.Ctx) error {

	var body subQuestionChoiceDto.CreateSubQuestionChoicesReq

	// ? Merge fiber http body with dto struct
	if err := c.BodyParser(&body); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, "The input data is invalid", nil)
	}

	// ? Validate field in body with dynamic function
	if err := validator.New().Struct(&body); err != nil {
		validationErrors := utils.TranslateValidationError(err.(validator.ValidationErrors))
		return response.ErrResponse(c, http.StatusBadRequest, "The input data is invalid", &validationErrors)
	}

	res, err := h.subQuestionChoiceUsecase.CreateSubQuestionChoice(&body)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}
