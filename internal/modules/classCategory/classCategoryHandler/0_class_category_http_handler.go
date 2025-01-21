package handler

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	classCategoryDto "github.com/gunktp20/digital-hubx-be/internal/modules/classCategory/classCategoryDto"
	classCategoryUsecase "github.com/gunktp20/digital-hubx-be/internal/modules/classCategory/classCategoryUsecase"
	"github.com/gunktp20/digital-hubx-be/pkg/response"
	"github.com/gunktp20/digital-hubx-be/pkg/utils"
)

type (
	ClassCategoryHttpHandlerService interface {
		CreateClassCategory(c *fiber.Ctx) error
		GetAllClassCategories(c *fiber.Ctx) error
		UpdateCategoryName(c *fiber.Ctx) error
	}

	classCategoryHttpHandler struct {
		classCategoryUsecase classCategoryUsecase.ClassCategoryUsecaseService
	}
)

func NewClassCategoryHttpHandler(usecase classCategoryUsecase.ClassCategoryUsecaseService) ClassCategoryHttpHandlerService {
	return &classCategoryHttpHandler{classCategoryUsecase: usecase}
}

// CreateClassCategory creates a new class category.
// @Summary Create a new class category
// @Description Allows an admin to create a new class category.
// @Tags Admin/Class Category
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "Create Class Category Request Body"
// @Success 200 {object} map[string]interface{} "Category created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /admin/class-category [post]
func (h *classCategoryHttpHandler) CreateClassCategory(c *fiber.Ctx) error {

	var body classCategoryDto.CreateClassCategoryReq

	// ? Merge fiber http body with dto struct
	if err := c.BodyParser(&body); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, "The input data is invalid", nil)
	}

	// ? Validate field in body with dynamic function
	if err := validator.New().Struct(&body); err != nil {
		validationErrors := utils.TranslateValidationError(err.(validator.ValidationErrors))
		return response.ErrResponse(c, http.StatusBadRequest, "The input data is invalid", &validationErrors)
	}

	res, err := h.classCategoryUsecase.CreateClassCategory(&body)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

// GetAllClassCategories retrieves all class categories.
// @Summary Get all class categories
// @Description Retrieves a paginated list of class categories.
// @Tags ClassCategory
// @Accept json
// @Produce json
// @Param keyword query string false "Keyword to filter categories"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Success 200 {object} map[string]interface{} "Paginated categories"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /class-category [get]
func (h *classCategoryHttpHandler) GetAllClassCategories(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	res, total, err := h.classCategoryUsecase.GetAllClassCategories(c.Query("keyword"), page, limit)
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

// UpdateCategoryName updates the name of an existing class category.
// @Summary Update class category name
// @Description Allows an admin to update the name of a class category.
// @Tags Admin/Class Category
// @Accept json
// @Produce json
// @Param category_id path string true "Category ID"
// @Param body body map[string]interface{} true "Update Category Name Request Body"
// @Success 200 {object} map[string]interface{} "Category name updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 404 {object} map[string]interface{} "Category not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /admin/class-category/{category_id} [put]
func (h *classCategoryHttpHandler) UpdateCategoryName(c *fiber.Ctx) error {

	var body classCategoryDto.UpdateCategoryNameReq
	if err := c.BodyParser(&body); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, "Invalid request body", nil)
	}

	// Validate input
	if err := validator.New().Struct(&body); err != nil {
		validationErrors := utils.TranslateValidationError(err.(validator.ValidationErrors))
		return response.ErrResponse(c, http.StatusBadRequest, "Validation failed", &validationErrors)
	}

	// Extract category ID from the URL
	categoryID := c.Params("category_id")

	err := h.classCategoryUsecase.UpdateCategoryName(categoryID, body.ClassCategoryName)
	if err != nil {
		if strings.Contains(err.Error(), "category name already exists") {
			return response.ErrResponse(c, http.StatusConflict, err.Error(), nil)
		}
		if strings.Contains(err.Error(), "category not found") {
			return response.ErrResponse(c, http.StatusNotFound, err.Error(), nil)
		}
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessResponse(c, http.StatusOK, fiber.Map{
		"message": "Category name updated successfully",
	})
}
