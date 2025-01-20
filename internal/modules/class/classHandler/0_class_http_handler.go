package handler

import (
	"math"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	classDto "github.com/gunktp20/digital-hubx-be/internal/modules/class/classDto"
	classUsecase "github.com/gunktp20/digital-hubx-be/internal/modules/class/classUsecase"
	"github.com/gunktp20/digital-hubx-be/pkg/models"
	"github.com/gunktp20/digital-hubx-be/pkg/response"
	"github.com/gunktp20/digital-hubx-be/pkg/utils"
)

var getContextAuth = utils.GetContextAuth

type (
	ClassHttpHandlerService interface {
		CreateClass(c *fiber.Ctx) error
		GetAllClasses(c *fiber.Ctx) error
		GetClassById(c *fiber.Ctx) error
		ToggleClassEnableQuestion(c *fiber.Ctx) error
		UpdateClassDetails(c *fiber.Ctx) error
		UpdateClassCoverImage(c *fiber.Ctx) error
		DeleteClass(c *fiber.Ctx) error
	}

	classHttpHandler struct {
		classUsecase classUsecase.ClassUsecaseService
	}
)

func NewClassHttpHandler(usecase classUsecase.ClassUsecaseService) ClassHttpHandlerService {
	return &classHttpHandler{classUsecase: usecase}
}

func (h *classHttpHandler) CreateClass(c *fiber.Ctx) error {

	fileHeader, err := c.FormFile("cover_image")
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, "Class cover image is required", nil)
	}

	// ? Convert Multipart file to bytes
	fileBytes, err := utils.ConvertMultipartFileToBytes(fileHeader)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, "Failed to convert multipart file to bytes", nil)
	}
	// ? Define allowed types & max file size
	allowedTypes := []string{"image/png", "image/jpg"}
	maxFileSize := int64(5 * 1024 * 1024)

	classLevelStr := c.FormValue("class_level")
	var classLevel int
	// ? If class level is sent strconv atoi for check class level is an integer
	if classLevelStr != "" {
		classLevel, err = strconv.Atoi(classLevelStr)
		if err != nil {
			return response.ErrResponse(c, http.StatusBadRequest, "Invalid class_level, must be an integer", nil)
		}
	}

	body := classDto.CreateClassReq{
		Title:           c.FormValue("title"),
		Description:     c.FormValue("description"),
		CoverImage:      "",
		ClassCategoryID: c.FormValue("class_category_id"),
		ClassLevel:      classLevel, // classLevel จะเป็น 0 ถ้าไม่ได้รับค่า
		ClassTier:       models.ClassTier(c.FormValue("class_tier")),
	}

	// ? Validate file with allowed types & max file size
	if err := utils.ValidateFile(fileBytes, allowedTypes, maxFileSize); err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	// ? Merge fiber http body with dto struct
	if err := c.BodyParser(&body); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, "The input data is invalid", nil)
	}

	// ? Validate field in body with dynamic function
	if err := validator.New().Struct(&body); err != nil {
		validationErrors := utils.TranslateValidationError(err.(validator.ValidationErrors))
		return response.ErrResponse(c, http.StatusBadRequest, "The input data is invalid", &validationErrors)
	}

	res, err := h.classUsecase.CreateClass(&body, fileBytes, fileHeader)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *classHttpHandler) GetAllClasses(c *fiber.Ctx) error {
	// อ่านค่า page และ limit
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	// อ่านค่า class_level และแปลงเป็น *int
	var classLevel *int
	if level := c.Query("class_level"); level != "" {
		parsedLevel, err := strconv.Atoi(level)
		if err != nil {
			return response.ErrResponse(c, http.StatusBadRequest, "Invalid class_level value", nil)
		}
		classLevel = &parsedLevel
	}

	// เรียกใช้ usecase
	res, total, err := h.classUsecase.GetAllClasses(
		c.Query("class_tier"),
		c.Query("keyword"),
		classLevel,
		c.Query("class_category"),
		page,
		limit,
	)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	// ส่งผลลัพธ์กลับ
	return response.SuccessResponse(c, http.StatusOK, fiber.Map{
		"data":       res,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": int(math.Ceil(float64(total) / float64(limit))),
	})
}

func (h *classHttpHandler) GetClassById(c *fiber.Ctx) error {

	res, err := h.classUsecase.GetClassById(c.Params("class_id"))
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *classHttpHandler) ToggleClassEnableQuestion(c *fiber.Ctx) error {

	newState, err := h.classUsecase.ToggleClassEnableQuestion(c.Params("class_id"))
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	var message string
	if newState {
		message = "EnableQuestion is now ENABLED"
	} else {
		message = "EnableQuestion is now DISABLED"
	}

	return response.SuccessResponse(c, http.StatusOK, &fiber.Map{
		"message": message,
	})
}

func (h *classHttpHandler) UpdateClassDetails(c *fiber.Ctx) error {
	var body classDto.UpdateClassReq

	if err := c.BodyParser(&body); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, "Invalid input data", nil)
	}

	classID := c.Params("class_id")
	if classID == "" {
		return response.ErrResponse(c, http.StatusBadRequest, "class ID is required", nil)
	}

	err := h.classUsecase.UpdateClassDetails(classID, body.Title, body.Description, body.ClassCategoryName, body.ClassTier, body.ClassLevel)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessResponse(c, http.StatusOK, &fiber.Map{
		"message": "Class details updated successfully",
	})
}

func (h *classHttpHandler) UpdateClassCoverImage(c *fiber.Ctx) error {

	fileHeader, err := c.FormFile("new_cover_image")
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, "New cover image is required", nil)
	}

	// ? Convert Multipart file to bytes
	fileBytes, err := utils.ConvertMultipartFileToBytes(fileHeader)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, "Failed to convert multipart file to bytes", nil)
	}
	// ? Define allowed types & max file size
	allowedTypes := []string{"image/png", "image/jpg", "image/svg"}
	maxFileSize := int64(5 * 1024 * 1024)

	// ? Validate file with allowed types & max file size
	if err := utils.ValidateFile(fileBytes, allowedTypes, maxFileSize); err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	err = h.classUsecase.UpdateClassCoverImage(c.Params("class_id"), fileBytes, fileHeader)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessResponse(c, http.StatusOK, &fiber.Map{
		"message": "Class cover image updated successfully",
	})
}

func (h *classHttpHandler) DeleteClass(c *fiber.Ctx) error {
	// ดึง class_id จาก URL
	classID := c.Params("class_id")
	if classID == "" {
		return response.ErrResponse(c, http.StatusBadRequest, "class ID is required", nil)
	}

	// เรียก Usecase เพื่อลบคลาส
	err := h.classUsecase.SoftDeleteClass(classID)
	if err != nil {
		if err.Error() == "class not found" {
			return response.ErrResponse(c, http.StatusNotFound, err.Error(), nil)
		}
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessResponse(c, http.StatusOK, fiber.Map{
		"message": "Class soft deleted successfully",
	})
}
