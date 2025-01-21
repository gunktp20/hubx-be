package dto

import "time"

type (
	// CreateClassCategoryReq represents the request body for creating a new class category.
	// @Description Request body for creating a class category.
	// @Param class_category_name body string true "Name of the class category"
	CreateClassCategoryReq struct {
		ClassCategoryName string `json:"class_category_name" validate:"required"`
	}

	// UpdateCategoryNameReq represents the request body for updating a class category name.
	// @Description Request body for updating a class category name.
	// @Param new_class_category_name body string true "New name of the class category"
	UpdateCategoryNameReq struct {
		ClassCategoryName string `json:"new_class_category_name" validate:"required"`
	}

	// CreateClassCategoryRes represents the response body after creating a new class category.
	// @Description Response body after creating a new class category.
	// @Param id query string true "Class category ID"
	// @Param class_category_name query string true "Name of the class category"
	// @Param created_at query string true "Category creation timestamp"
	// @Param updated_at query string true "Category last updated timestamp"
	CreateClassCategoryRes struct {
		ID                string    `json:"id"`
		ClassCategoryName string    `json:"class_category_name"`
		CreatedAt         time.Time `json:"created_at"`
		UpdatedAt         time.Time `json:"updated_at"`
	}
)
