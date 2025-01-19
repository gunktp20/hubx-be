package handler

import (
	"math"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	subQuestionDto "github.com/gunktp20/digital-hubx-be/internal/modules/subQuestion/subQuestionDto"
	subQuestionUsecase "github.com/gunktp20/digital-hubx-be/internal/modules/subQuestion/subQuestionUsecase"
	"github.com/gunktp20/digital-hubx-be/pkg/response"
	"github.com/gunktp20/digital-hubx-be/pkg/utils"
)

type (
	SubQuestionHttpHandlerService interface {
		CreateSubQuestion(c *fiber.Ctx) error
		GetSubQuestionsByQuestionID(c *fiber.Ctx) error
		GetSubQuestionsByChoiceID(c *fiber.Ctx) error
	}

	subQuestionHttpHandler struct {
		subQuestionUsecase subQuestionUsecase.SubQuestionUsecaseService
	}
)

func NewSubQuestionHttpHandler(usecase subQuestionUsecase.SubQuestionUsecaseService) SubQuestionHttpHandlerService {
	return &subQuestionHttpHandler{subQuestionUsecase: usecase}
}

func (h *subQuestionHttpHandler) CreateSubQuestion(c *fiber.Ctx) error {

	var body subQuestionDto.CreateSubQuestionReq

	// ? Merge fiber http body with dto struct
	if err := c.BodyParser(&body); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, "The input data is invalid", nil)
	}

	// ? Validate field in body with dynamic function
	if err := validator.New().Struct(&body); err != nil {
		validationErrors := utils.TranslateValidationError(err.(validator.ValidationErrors))
		return response.ErrResponse(c, http.StatusBadRequest, "The input data is invalid", &validationErrors)
	}

	res, err := h.subQuestionUsecase.CreateSubQuestion(&body)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *subQuestionHttpHandler) GetSubQuestionsByQuestionID(c *fiber.Ctx) error {

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	res, total, err := h.subQuestionUsecase.GetSubQuestionsByQuestionID(c.Params("question_id"), page, limit)

	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessResponse(c, http.StatusOK, fiber.Map{
		"data":       res,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": int(math.Ceil(float64(total) / float64(limit))),
	})
}

func (h *subQuestionHttpHandler) GetSubQuestionsByChoiceID(c *fiber.Ctx) error {

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	res, total, err := h.subQuestionUsecase.GetSubQuestionsByChoiceID(c.Params("choice_id"), page, limit)

	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessResponse(c, http.StatusOK, fiber.Map{
		"data":       res,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": int(math.Ceil(float64(total) / float64(limit))),
	})
}
