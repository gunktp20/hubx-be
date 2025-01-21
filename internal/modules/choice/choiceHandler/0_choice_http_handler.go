package handler

import (
	"math"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	choiceDto "github.com/gunktp20/digital-hubx-be/internal/modules/choice/choiceDto"
	choiceUsecase "github.com/gunktp20/digital-hubx-be/internal/modules/choice/choiceUsecase"
	"github.com/gunktp20/digital-hubx-be/pkg/response"
	"github.com/gunktp20/digital-hubx-be/pkg/utils"
)

type (
	ChoiceHttpHandlerService interface {
		CreateChoice(c *fiber.Ctx) error
	}

	choiceHttpHandler struct {
		choiceUsecase choiceUsecase.ChoiceUsecaseService
	}
)

func NewChoiceHttpHandler(usecase choiceUsecase.ChoiceUsecaseService) ChoiceHttpHandlerService {
	return &choiceHttpHandler{choiceUsecase: usecase}
}

// CreateChoice creates a choice for a specific question.
// @Summary Create a new choice
// @Description Allows an admin to create a choice for a specific question.
// @Tags Admin/Choice
// @Accept json
// @Produce json
// @Param body body choiceDto.CreateChoiceReq true "Create Choice Request Body"
// @Success 200 {object} map[string]interface{} "Operation successful"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /admin/choice [post]
func (h *choiceHttpHandler) CreateChoice(c *fiber.Ctx) error {

	var body choiceDto.CreateChoiceReq

	// ? Merge fiber http body with dto struct
	if err := c.BodyParser(&body); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, "The input data is invalid", nil)
	}

	// ? Validate field in body with dynamic function
	if err := validator.New().Struct(&body); err != nil {
		validationErrors := utils.TranslateValidationError(err.(validator.ValidationErrors))
		return response.ErrResponse(c, http.StatusBadRequest, "The input data is invalid", &validationErrors)
	}

	res, err := h.choiceUsecase.CreateChoice(&body)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *choiceHttpHandler) GetChoicesByClassID(c *fiber.Ctx) error {

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	res, total, err := h.choiceUsecase.GetChoicesByClassID(c.Params("class_id"), page, limit)

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
