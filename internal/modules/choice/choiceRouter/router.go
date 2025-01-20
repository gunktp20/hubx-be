package router

import (
	"github.com/gofiber/fiber/v2"
	choiceHandler "github.com/gunktp20/digital-hubx-be/internal/modules/choice/choiceHandler"
)

func SetChoiceRoutes(api fiber.Router, choiceHttpHandler choiceHandler.ChoiceHttpHandlerService) {
	_ = api.Group("/choice")

	// ? Admin Routes Group
	adminRoute := api.Group("/admin/choice")
	adminRoute.Post("/", choiceHttpHandler.CreateChoice)
}
