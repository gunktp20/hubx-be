package router

import (
	"github.com/gofiber/fiber/v2"
	choiceHandler "github.com/gunktp20/digital-hubx-be/internal/modules/choice/choiceHandler"
	"github.com/gunktp20/digital-hubx-be/pkg/middleware"
)

func SetChoiceRoutes(api fiber.Router, choiceHttpHandler choiceHandler.ChoiceHttpHandlerService) {
	_ = api.Group("/choice")

	// ? Admin Routes
	adminRoute := api.Group("/admin/choice", middleware.Ident, middleware.PermissionCheck)
	adminRoute.Post("/", choiceHttpHandler.CreateChoice)
}
