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
