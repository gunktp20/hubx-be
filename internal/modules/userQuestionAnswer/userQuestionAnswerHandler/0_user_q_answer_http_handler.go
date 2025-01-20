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

// func (h *userQuestionAnswerHttpHandler) CreateUserQuestionAnswer(c *fiber.Ctx) error {

// 	userIDStr, err := utils.GetUserIDFromContext(c)
// 	if err != nil {
// 		return response.ErrResponse(c, http.StatusUnauthorized, err.Error(), nil)
// 	}

// 	var body userQuestionAnswerDto.CreateUserQuestionAnswerReq

// 	// ? Merge fiber http body with dto struct
// 	if err := c.BodyParser(&body); err != nil {
// 		return response.ErrResponse(c, http.StatusBadRequest, "The input data is invalid", nil)
// 	}

// 	// ? Validate field in body with dynamic function
// 	if err := validator.New().Struct(&body); err != nil {
// 		validationErrors := utils.TranslateValidationError(err.(validator.ValidationErrors))
// 		return response.ErrResponse(c, http.StatusBadRequest, "The input data is invalid", &validationErrors)
// 	}

// 	res, err := h.userQuestionAnswerUsecase.CreateUserQuestionAnswer(&body, userIDStr)
// 	if err != nil {
// 		return response.ErrResponse(c, http.StatusInternalServerError, err.Error(), nil)
// 	}

// 	return response.SuccessResponse(c, http.StatusOK, res)
// }

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
