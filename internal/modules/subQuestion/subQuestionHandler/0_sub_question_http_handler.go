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

// @Summary Create a new sub-question
// @Description Allows an admin to create a sub-question for a specific choice.
// @Tags Admin/SubQuestion
// @Accept json
// @Produce json
// @Param body body subQuestionDto.CreateSubQuestionReq true "Create SubQuestion Request Body"
// @Success 200 {object} subQuestionDto.CreateSubQuestionRes "Sub-question created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input" example:{"message":"Invalid input","status":400,"details":{"field":"error description"}}
// @Failure 500 {object} map[string]interface{} "Internal Server Error" example:{"message":"Internal Server Error","status":500,"details":null}
// @Router /admin/sub-question [post]
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

// @Summary Get sub-questions by question ID
// @Description Retrieves sub-questions for a given question ID.
// @Tags SubQuestion
// @Accept json
// @Produce json
// @Param question_id path string true "Question ID" format(uuid)
// @Param page query int false "Page number" example:1
// @Param limit query int false "Number of items per page" example:10
// @Success 200 {object} map[string]interface{} "List of sub-questions" example:{"data":[],"total":0,"page":1,"limit":10,"totalPages":1}
// @Failure 500 {object} map[string]interface{} "Internal Server Error" example:{"message":"Internal Server Error","status":500,"details":null}
// @Router /sub-question/{question_id}/question [get]
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

// @Summary Get sub-questions by choice ID
// @Description Retrieves sub-questions for a given choice ID.
// @Tags SubQuestion
// @Accept json
// @Produce json
// @Param choice_id path string true "Choice ID" format(uuid)
// @Param page query int false "Page number" example:1
// @Param limit query int false "Number of items per page" example:10
// @Success 200 {object} map[string]interface{} "List of sub-questions" example:{"data":[],"total":0,"page":1,"limit":10,"totalPages":1}
// @Failure 500 {object} map[string]interface{} "Internal Server Error" example:{"message":"Internal Server Error","status":500,"details":null}
// @Router /sub-question/{choice_id}/choice [get]
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
