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

// CreateClass creates a new class.
// @Summary Create a new class
// @Description Allows an admin to create a new class by providing title, description, category, tier, and other details.
// @Tags Admin/Class
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Class Title"
// @Param description formData string true "Class Description"
// @Param cover_image formData file true "Class Cover Image"
// @Param class_category_id formData string true "Class Category ID"
// @Param class_tier formData string true "Class Tier" enum(Essential,Literacy,Mastery)
// @Param class_level formData int false "Class Level"
// @Success 200 {object} map[string]interface{} "Operation successful"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /admin/class [post]
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

// GetAllClasses retrieves a paginated list of classes.
// @Summary Get all classes
// @Description Fetch all classes with optional filters like tier, keyword, and category.
// @Tags Class
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10)"
// @Param class_tier query string false "Class Tier" enum(Essential,Literacy,Mastery)
// @Param keyword query string false "Search keyword"
// @Param class_level query int false "Filter by class level"
// @Param class_category query string false "Filter by class category"
// @Router /class [get]
func (h *classHttpHandler) GetAllClasses(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	var classLevel *int
	if level := c.Query("class_level"); level != "" {
		parsedLevel, err := strconv.Atoi(level)
		if err != nil {
			return response.ErrResponse(c, http.StatusBadRequest, "Invalid class_level value", nil)
		}
		classLevel = &parsedLevel
	}

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

	return response.SuccessResponse(c, http.StatusOK, fiber.Map{
		"data":       res,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": int(math.Ceil(float64(total) / float64(limit))),
	})
}

// GetClassById retrieves a class by its ID.
// @Summary Get class by ID
// @Description Fetch the details of a specific class by its ID.
// @Tags Class
// @Accept json
// @Produce json
// @Param class_id path string true "Class ID"
// @Success 200 {object} classDto.CreateClassRes "Class details"
// @Failure 404 {object} response.MsgResponse "Class not found"
// @Failure 500 {object} response.MsgResponse "Internal Server Error"
// @Router /class/{class_id} [get]
func (h *classHttpHandler) GetClassById(c *fiber.Ctx) error {

	res, err := h.classUsecase.GetClassById(c.Params("class_id"))
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

// ToggleClassEnableQuestion toggles the question enablement status of a class.
// @Summary Toggle EnableQuestion status
// @Description Enables or disables question functionality for a specific class.
// @Tags Class
// @Accept json
// @Produce json
// @Param class_id path string true "Class ID"
// @Success 200 {object} map[string]interface{} "Operation successful"
// @Failure 404 {object} map[string]interface{} "Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /admin/class/{class_id}/toggle-enable-question [put]
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

// CreateClassCategory creates a new class category.
// @Summary Create a new class category
// @Description Allows an admin to create a new class category.
// @Tags Admin/Class
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "Create Class Category Request Body"
// @Success 200 {object} map[string]interface{} "Category created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /admin/class-category [post]
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

// UpdateClassCoverImage updates the cover image of a class.
// @Summary Update class cover image
// @Description Allows an admin to update the cover image of a class.
// @Tags Class
// @Accept multipart/form-data
// @Produce json
// @Param class_id path string true "Class ID"
// @Param new_cover_image formData file true "New Cover Image"
// @Success 200 {object} map[string]interface{} "Operation successful"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 404 {object} map[string]interface{} "Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /admin/class/{class_id}/cover-image [put]
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

// DeleteClass soft deletes a class by its ID.
// @Summary Delete a class
// @Description Soft deletes a class by marking it as removed.
// @Tags Admin/Class
// @Accept json
// @Produce json
// @Param class_id path string true "Class ID"
// @Success 200 {object} map[string]interface{} "Operation successful"
// @Failure 404 {object} map[string]interface{} "Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /admin/class/{class_id} [delete]
func (h *classHttpHandler) DeleteClass(c *fiber.Ctx) error {

	classID := c.Params("class_id")
	if classID == "" {
		return response.ErrResponse(c, http.StatusBadRequest, "class ID is required", nil)
	}

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
