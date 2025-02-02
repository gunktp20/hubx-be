package router

import (
	"github.com/gofiber/fiber/v2"
	classCategoryHandler "github.com/gunktp20/digital-hubx-be/internal/modules/classCategory/classCategoryHandler"
	"github.com/gunktp20/digital-hubx-be/pkg/middleware"
)

func SetClassCategoryRoutes(api fiber.Router, classCategoryHttpHandler classCategoryHandler.ClassCategoryHttpHandlerService) {
	classCategoryRoute := api.Group("/class-category")

	classCategoryRoute.Get("/", classCategoryHttpHandler.GetAllClassCategories)

	// ? Admin Routes Group
	adminRoute := api.Group("/admin/class-category", middleware.Ident, middleware.PermissionCheck)
	adminRoute.Post("/", classCategoryHttpHandler.CreateClassCategory)
	adminRoute.Put("/:category_id", classCategoryHttpHandler.UpdateCategoryName)

}
