package handler

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	classSessionDto "github.com/gunktp20/digital-hubx-be/internal/modules/classSession/classSessionDto"
	classSessionUsecase "github.com/gunktp20/digital-hubx-be/internal/modules/classSession/classSessionUsecase"
	"github.com/gunktp20/digital-hubx-be/pkg/response"
	"github.com/gunktp20/digital-hubx-be/pkg/utils"
)

type (
	ClassSessionHttpHandlerService interface {
		CreateClassSession(c *fiber.Ctx) error
		GetAllClassSessions(c *fiber.Ctx) error
		SetMaxCapacity(c *fiber.Ctx) error
		UpdateClassSessionLocation(c *fiber.Ctx) error
	}

	classSessionHttpHandler struct {
		classSessionUsecase classSessionUsecase.ClassSessionUsecaseService
	}
)

func NewClassSessionHttpHandler(usecase classSessionUsecase.ClassSessionUsecaseService) ClassSessionHttpHandlerService {
	return &classSessionHttpHandler{classSessionUsecase: usecase}
}

func (h *classSessionHttpHandler) CreateClassSession(c *fiber.Ctx) error {

	var body classSessionDto.CreateClassSessionReq

	// ? Merge fiber http body with dto struct
	if err := c.BodyParser(&body); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, "The input data is invalid", nil)
	}

	// ? Validate field in body with dynamic function
	if err := validator.New().Struct(&body); err != nil {
		validationErrors := utils.TranslateValidationError(err.(validator.ValidationErrors))
		return response.ErrResponse(c, http.StatusBadRequest, "The input data is invalid", &validationErrors)
	}

	res, err := h.classSessionUsecase.CreateClassSession(&body)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *classSessionHttpHandler) GetAllClassSessions(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	res, total, err := h.classSessionUsecase.GetAllClassSessions(c.Query("class_id"), c.Query("class_tier"), page, limit)
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

func (h *classSessionHttpHandler) SetMaxCapacity(c *fiber.Ctx) error {

	var body classSessionDto.SetMaxCapacityReq
	// ? Merge fiber http body with dto struct
	if err := c.BodyParser(&body); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, "The input data is invalid", nil)
	}

	// ? Validate field in body with dynamic function
	if err := validator.New().Struct(&body); err != nil {
		validationErrors := utils.TranslateValidationError(err.(validator.ValidationErrors))
		return response.ErrResponse(c, http.StatusBadRequest, "The input data is invalid", &validationErrors)
	}

	err := h.classSessionUsecase.SetMaxCapacity(c.Params("class_session_id"), body.NewCapacity)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessResponse(c, http.StatusOK, &fiber.Map{
		"message": "The max capacity was updated successfully",
	})
}

func (h *classSessionHttpHandler) UpdateClassSessionLocation(c *fiber.Ctx) error {
	classSessionID := c.Params("class_session_id")

	// Parse the request body
	var body struct {
		NewLocation string `json:"new_location"`
	}

	if err := c.BodyParser(&body); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, "Invalid request body", nil)
	}

	// Validate input
	if body.NewLocation == "" {
		return response.ErrResponse(c, http.StatusBadRequest, "New location is required", nil)
	}

	// Call use case to update location
	err := h.classSessionUsecase.UpdateLocation(classSessionID, body.NewLocation)
	if err != nil {
		if strings.Contains(err.Error(), "no class session found") {
			return response.ErrResponse(c, http.StatusNotFound, err.Error(), nil)
		}
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessResponse(c, http.StatusOK, fiber.Map{
		"message": "Class session location updated successfully",
	})
}
