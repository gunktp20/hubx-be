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
		DeleteClassSessionByID(c *fiber.Ctx) error
	}

	classSessionHttpHandler struct {
		classSessionUsecase classSessionUsecase.ClassSessionUsecaseService
	}
)

func NewClassSessionHttpHandler(usecase classSessionUsecase.ClassSessionUsecaseService) ClassSessionHttpHandlerService {
	return &classSessionHttpHandler{classSessionUsecase: usecase}
}

// @Summary Create a new class session
// @Description Allows an admin to create a new class session for a class.
// @Tags Admin/Class Session
// @Accept json
// @Produce json
// @Param body body classSessionDto.CreateClassSessionReq true "Create Class Session Request Body"
// @Success 200 {object} classSessionDto.CreateClassSessionRes "Class session created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input" example:{"message":"Invalid input","status":400,"details":{"field":"error description"}}
// @Failure 500 {object} map[string]interface{} "Internal Server Error" example:{"message":"Internal Server Error","status":500,"details":null}
// @Router /admin/class-session [post]
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

// @Summary Get all class sessions
// @Description Retrieves a list of class sessions with optional filters.
// @Tags Class Session
// @Accept json
// @Produce json
// @Param class_id query string false "Class ID"
// @Param class_tier query string false "Class tier"
// @Param page query int false "Page number" example:1
// @Param limit query int false "Number of items per page" example:10
// @Success 200 {object} map[string]interface{} "List of class sessions" example:{"data":[],"total":0,"page":1,"limit":10,"totalPages":1}
// @Failure 500 {object} map[string]interface{} "Internal Server Error" example:{"message":"Internal Server Error","status":500,"details":null}
// @Router /class-session [get]
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

// @Summary Update max capacity for a class session
// @Description Allows an admin to update the maximum capacity for a class session.
// @Tags Admin/Class Session
// @Accept json
// @Produce json
// @Param class_session_id path string true "Class session ID"
// @Param body body classSessionDto.SetMaxCapacityReq true "Max Capacity Update Request Body"
// @Success 200 {object} map[string]interface{} "Max capacity updated successfully" example:{"message":"The max capacity was updated successfully","status":200,"details":null}
// @Failure 400 {object} map[string]interface{} "Invalid input" example:{"message":"Invalid input","status":400,"details":{"field":"error description"}}
// @Failure 500 {object} map[string]interface{} "Internal Server Error" example:{"message":"Internal Server Error","status":500,"details":null}
// @Router /admin/class-session/{class_session_id}/max-capacity [put]
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

// @Summary Update location for a class session
// @Description Allows an admin to update the location for a class session.
// @Tags Admin/Class Session
// @Accept json
// @Produce json
// @Param class_session_id path string true "Class session ID"
// @Param body body classSessionDto.UpdateClassSessionLocation true "Location Update Request Body"
// @Success 200 {object} map[string]interface{} "Location updated successfully" example:{"message":"Class session location updated successfully","status":200,"details":null}
// @Failure 400 {object} map[string]interface{} "Invalid input" example:{"message":"Invalid input","status":400,"details":{"field":"error description"}}
// @Failure 500 {object} map[string]interface{} "Internal Server Error" example:{"message":"Internal Server Error","status":500,"details":null}
// @Router /admin/class-session/{class_session_id}/location [put]
func (h *classSessionHttpHandler) UpdateClassSessionLocation(c *fiber.Ctx) error {
	classSessionID := c.Params("class_session_id")

	// Parse the request body
	var body classSessionDto.UpdateClassSessionLocation
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

// @Summary Delete a class session
// @Description Allows an admin to delete a specific class session by ID.
// @Tags Admin/Class Session
// @Accept json
// @Produce json
// @Param class_session_id path string true "Class session ID"
// @Success 200 {object} map[string]interface{} "Class session deleted successfully" example:{"message":"Class session deleted successfully","status":200,"details":null}
// @Failure 500 {object} map[string]interface{} "Internal Server Error" example:{"message":"Internal Server Error","status":500,"details":null}
// @Router /admin/class-session/{class_session_id} [delete]
func (h *classSessionHttpHandler) DeleteClassSessionByID(c *fiber.Ctx) error {
	classSessionID := c.Params("class_session_id")

	if classSessionID == "" {
		return response.ErrResponse(c, fiber.StatusBadRequest, "Class session ID is required", nil)
	}

	err := h.classSessionUsecase.DeleteClassSessionByID(classSessionID)
	if err != nil {
		return response.ErrResponse(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Class session deleted successfully",
	})
}
