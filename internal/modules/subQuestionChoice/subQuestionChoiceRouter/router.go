package router

import (
	"github.com/gofiber/fiber/v2"
	subQuestionChoiceHandler "github.com/gunktp20/digital-hubx-be/internal/modules/subQuestionChoice/subQuestionChoiceHandler"
)

func SetSubQuestionChoiceRoutes(api fiber.Router, subQuestionChoiceHttpHandler subQuestionChoiceHandler.SubQuestionChoiceHttpHandlerService) {
	_ = api.Group("/sub-question-choice")

	// ? Admin Routes Group
	adminRoute := api.Group("/admin/class-session")
	adminRoute.Post("/", subQuestionChoiceHttpHandler.CreateSubQuestionChoice)

}
