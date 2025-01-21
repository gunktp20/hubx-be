package handler

import (
	"math"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	attendanceDto "github.com/gunktp20/digital-hubx-be/internal/modules/attendance/attendanceDto"
	attendanceUsecase "github.com/gunktp20/digital-hubx-be/internal/modules/attendance/attendanceUsecase"
	"github.com/gunktp20/digital-hubx-be/pkg/response"
	"github.com/gunktp20/digital-hubx-be/pkg/utils"
)

type (
	AttendanceHttpHandlerService interface {
		CreateAttendance(c *fiber.Ctx) error
	}

	attendanceHttpHandler struct {
		attendanceUsecase attendanceUsecase.AttendanceUsecaseService
	}
)

func NewAttendanceHttpHandler(usecase attendanceUsecase.AttendanceUsecaseService) AttendanceHttpHandlerService {
	return &attendanceHttpHandler{attendanceUsecase: usecase}
}

// CreateAttendance creates an attendance record for a class session.
// @Summary Create a new attendance record
// @Description Allows an admin to create an attendance record for a specific class session.
// @Tags Admin/Attendance
// @Accept json
// @Produce json
// @Param body body attendanceDto.CreateAttendanceReq true "Create Attendance Request Body"
// @Success 200 {object} map[string]interface{} "Operation successful"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /admin/attendance [post]
func (h *attendanceHttpHandler) CreateAttendance(c *fiber.Ctx) error {
	var body attendanceDto.CreateAttendanceReq

	// ? Merge fiber http body with dto struct
	if err := c.BodyParser(&body); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, "The input data is invalid", nil)
	}

	// ? Validate field in body with dynamic function
	if err := validator.New().Struct(&body); err != nil {
		validationErrors := utils.TranslateValidationError(err.(validator.ValidationErrors))
		return response.ErrResponse(c, http.StatusBadRequest, "The input data is invalid", &validationErrors)
	}

	res, err := h.attendanceUsecase.CreateAttendance(&body)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *attendanceHttpHandler) GetAttendancesByClassID(c *fiber.Ctx) error {

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	res, total, err := h.attendanceUsecase.GetAttendancesByClassID(c.Params("class_id"), page, limit)

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
