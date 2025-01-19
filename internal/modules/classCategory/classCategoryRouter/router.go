package router

import (
	"github.com/gofiber/fiber/v2"
	classCategoryHandler "github.com/gunktp20/digital-hubx-be/internal/modules/classCategory/classCategoryHandler"
)

func SetClassCategoryRoutes(api fiber.Router, classCategoryHttpHandler classCategoryHandler.ClassCategoryHttpHandlerService) {
	routes := api.Group("/class-category")

	routes.Get("/", classCategoryHttpHandler.GetAllClassCategories)
	routes.Post("/", classCategoryHttpHandler.CreateClassCategory)

}
