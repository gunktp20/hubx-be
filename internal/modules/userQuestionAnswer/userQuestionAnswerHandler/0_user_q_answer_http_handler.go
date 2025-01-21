package handler

import (
	"math"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	userQuestionAnswerDto "github.com/gunktp20/digital-hubx-be/internal/modules/userQuestionAnswer/userQuestionAnswerDto"
	userQuestionAnswerUsecase "github.com/gunktp20/digital-hubx-be/internal/modules/userQuestionAnswer/userQuestionAnswerUsecase"
	"github.com/gunktp20/digital-hubx-be/pkg/response"
	"github.com/gunktp20/digital-hubx-be/pkg/utils"
)

var getContextAuth = utils.GetContextAuth

type (
	UserQuestionAnswerHttpHandlerService interface {
		// CreateUserQuestionAnswer(c *fiber.Ctx) error
		GetUserQuestionAnswersWithClassId(c *fiber.Ctx) error
		CreateMultipleUserQuestionAnswers(c *fiber.Ctx) error
	}

	userQuestionAnswerHttpHandler struct {
		userQuestionAnswerUsecase userQuestionAnswerUsecase.UserQuestionAnswerUsecaseService
	}
)

func NewUserQuestionAnswerHttpHandler(usecase userQuestionAnswerUsecase.UserQuestionAnswerUsecaseService) UserQuestionAnswerHttpHandlerService {
	return &userQuestionAnswerHttpHandler{userQuestionAnswerUsecase: usecase}
}

// @Summary Retrieve user question answers by class ID
// @Description Fetches a paginated list of user question answers for a specific class.
// @Tags UserQuestionAnswer
// @Accept json
// @Produce json
// @Param class_id path string true "Class ID"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10)"
// @Success 200 {object} map[string]interface{} "Success response" example:{"data":[...],"total":100,"page":1,"limit":10,"totalPages":10}
// @Failure 500 {object} map[string]interface{} "Internal Server Error" example:{"message":"Internal Server Error","status":500,"details":null}
// @Router /user-question-answer/{class_id}/class [get]
func (h *userQuestionAnswerHttpHandler) GetUserQuestionAnswersWithClassId(c *fiber.Ctx) error {
	_, _, userEmail := getContextAuth(c.UserContext())

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	res, total, err := h.userQuestionAnswerUsecase.GetUserQuestionAnswersWithClassId(userEmail, c.Params("class_id"), page, limit)
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

// @Summary Create multiple user question answers
// @Description Allows an admin to submit multiple question answers for a specific class.
// @Tags Admin/User Question Answer
// @Accept json
// @Produce json
// @Param class_id path string true "Class ID"
// @Param body body []userQuestionAnswerDto.CreateUserQuestionAnswerReq true "List of user question answers"
// @Success 200 {object} map[string]interface{} "Success response" example:{"message":"User question answers created successfully","status":200,"details":null}
// @Failure 400 {object} map[string]interface{} "Invalid input" example:{"message":"Invalid input","status":400,"details":{"field":"error description"}}
// @Failure 500 {object} map[string]interface{} "Internal Server Error" example:{"message":"Internal Server Error","status":500,"details":null}
// @Router /admin/user-question-answer/{class_id}/class [post]
func (h *userQuestionAnswerHttpHandler) CreateMultipleUserQuestionAnswers(c *fiber.Ctx) error {
	_, _, userEmail := getContextAuth(c.UserContext())

	var requestBody struct {
		Answers []userQuestionAnswerDto.CreateUserQuestionAnswerReq `json:"answers"`
	}

	// Parse JSON body
	if err := c.BodyParser(&requestBody); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, "The input data is invalid", nil)
	}

	// Validate each answer in the list
	for _, item := range requestBody.Answers {
		if err := validator.New().Struct(&item); err != nil {
			validationErrors := utils.TranslateValidationError(err.(validator.ValidationErrors))
			return response.ErrResponse(c, http.StatusBadRequest, "The input data is invalid", &validationErrors)
		}
	}

	// Call the use case with the whole list of answers
	res, err := h.userQuestionAnswerUsecase.CreateMultipleUserQuestionAnswers(
		requestBody.Answers, c.Params("class_id"), userEmail,
	)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}
