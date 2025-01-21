package router

import (
	"github.com/gofiber/fiber/v2"
	subQuestionChoiceHandler "github.com/gunktp20/digital-hubx-be/internal/modules/subQuestionChoice/subQuestionChoiceHandler"
	"github.com/gunktp20/digital-hubx-be/pkg/middleware"
)

func SetSubQuestionChoiceRoutes(api fiber.Router, subQuestionChoiceHttpHandler subQuestionChoiceHandler.SubQuestionChoiceHttpHandlerService) {
	_ = api.Group("/sub-question-choice")

	// ? Admin Routes Group
	adminRoute := api.Group("/admin/class-session", middleware.Ident, middleware.PermissionCheck)
	adminRoute.Post("/", subQuestionChoiceHttpHandler.CreateSubQuestionChoice)

}
