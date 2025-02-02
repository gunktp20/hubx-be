package handler

import (
	"math"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	classRegistrationDto "github.com/gunktp20/digital-hubx-be/internal/modules/classRegistration/classRegistrationDto"
	classRegistrationUsecase "github.com/gunktp20/digital-hubx-be/internal/modules/classRegistration/classRegistrationUsecase"
	"github.com/gunktp20/digital-hubx-be/pkg/response"
	"github.com/gunktp20/digital-hubx-be/pkg/utils"
)

var getContextAuth = utils.GetContextAuth

type (
	ClassRegistrationHttpHandlerService interface {
		CreateClassRegistration(c *fiber.Ctx) error
		GetUserRegistrations(c *fiber.Ctx) error
		CancelClassRegistration(c *fiber.Ctx) error
		ResetCancelledQuota(c *fiber.Ctx) error
		DeleteUserClassRegistrationBySession(c *fiber.Ctx) error
	}

	classRegistrationHttpHandler struct {
		classRegistrationUsecase classRegistrationUsecase.ClassRegistrationUsecaseService
	}
)

func NewClassRegistrationHttpHandler(usecase classRegistrationUsecase.ClassRegistrationUsecaseService) ClassRegistrationHttpHandlerService {
	return &classRegistrationHttpHandler{classRegistrationUsecase: usecase}
}

// @Summary Create a class registration
// @Description Allows a user to register for a class session.
// @Tags ClassRegistration
// @Accept json
// @Produce json
// @Param body body classRegistrationDto.CreateClassRegistrationReq true "Create ClassRegistration Request Body"
// @Success 200 {object} classRegistrationDto.CreateClassRegistrationRes "Registration created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /class-registration [post]
func (h *classRegistrationHttpHandler) CreateClassRegistration(c *fiber.Ctx) error {

	_, _, userEmail := getContextAuth(c.UserContext())
	var body classRegistrationDto.CreateClassRegistrationReq

	// ? Merge fiber http body with dto struct
	if err := c.BodyParser(&body); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, "The input data is invalid", nil)
	}

	// ? Validate field in body with dynamic function
	if err := validator.New().Struct(&body); err != nil {
		validationErrors := utils.TranslateValidationError(err.(validator.ValidationErrors))
		return response.ErrResponse(c, http.StatusBadRequest, "The input data is invalid", &validationErrors)
	}

	res, err := h.classRegistrationUsecase.CreateClassRegistration(&body, userEmail)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

// @Summary Get user registrations
// @Description Fetch paginated user registrations.
// @Tags ClassRegistration
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} map[string]interface{} "List of user registrations"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /class-registration [get]
func (h *classRegistrationHttpHandler) GetUserRegistrations(c *fiber.Ctx) error {
	_, _, userEmail := getContextAuth(c.UserContext())

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	res, total, err := h.classRegistrationUsecase.GetUserRegistrations(userEmail, page, limit)
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

// @Summary Cancel a class registration
// @Description Allows a user to cancel their class session registration.
// @Tags ClassRegistration
// @Accept json
// @Produce json
// @Param class_id path string true "Class ID"
// @Success 200 {object} map[string]interface{} "Registration cancelled successfully"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /class-registration/{class_id}/cancel [delete]
func (h *classRegistrationHttpHandler) CancelClassRegistration(c *fiber.Ctx) error {
	_, _, userEmail := getContextAuth(c.UserContext())

	classID := c.Params("class_id")
	if classID == "" {
		return response.ErrResponse(c, http.StatusBadRequest, "class_id is required", nil)
	}
	err := h.classRegistrationUsecase.CancelClassRegistration(userEmail, classID)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessResponse(c, http.StatusOK, &fiber.Map{
		"message": "Class registration cancelled successfully",
	})
}

// @Summary Reset user's cancellation quota
// @Description Allows an admin to reset a user's cancellation quota for a class.
// @Tags Admin/ClassRegistration
// @Accept json
// @Produce json
// @Param body body classRegistrationDto.ResetCancelledQuotaReq true "Reset Cancelled Quota Request Body"
// @Success 200 {object} map[string]interface{} "Cancellation quota reset successfully"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /admin/class-registration/reset-cancel-quota [post]
func (h *classRegistrationHttpHandler) ResetCancelledQuota(c *fiber.Ctx) error {

	var body classRegistrationDto.ResetCancelledQuotaReq

	// ? Merge fiber http body with dto struct
	if err := c.BodyParser(&body); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, "The input data is invalid", nil)
	}

	// ? Validate field in body with dynamic function
	if err := validator.New().Struct(&body); err != nil {
		validationErrors := utils.TranslateValidationError(err.(validator.ValidationErrors))
		return response.ErrResponse(c, http.StatusBadRequest, "The input data is invalid", &validationErrors)
	}

	err := h.classRegistrationUsecase.ResetCancelledQuota(&body)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessResponse(c, http.StatusOK, &fiber.Map{
		"message": "The user's cancellation quota for the specified class has been reset",
	})
}

// @Summary Delete a user's class registration
// @Description Allows an admin to delete a user's registration for a specific class session.
// @Tags Admin/ClassRegistration
// @Accept json
// @Produce json
// @Param class_session_id path string true "Class session ID"
// @Param email path string true "User's email"
// @Success 200 {object} map[string]interface{} "Registration deleted successfully"
// @Failure 404 {object} map[string]interface{} "Registration not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /admin/class-registration/{class_session_id}/{email} [delete]
func (h *classRegistrationHttpHandler) DeleteUserClassRegistrationBySession(c *fiber.Ctx) error {
	userEmail := c.Query("user_email")
	classSessionID := c.Query("class_session_id")

	if userEmail == "" || classSessionID == "" {
		return response.ErrResponse(c, http.StatusBadRequest, "user_email and class_session_id are required", nil)
	}

	err := h.classRegistrationUsecase.DeleteUserClassRegistrationBySession(userEmail, classSessionID)
	if err != nil {
		if err.Error() == "no user class registration found for the provided user email and class session ID" {
			return response.ErrResponse(c, http.StatusNotFound, err.Error(), nil)
		}
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessResponse(c, http.StatusOK, &fiber.Map{
		"message": "User class registration deleted successfully",
	})
}
