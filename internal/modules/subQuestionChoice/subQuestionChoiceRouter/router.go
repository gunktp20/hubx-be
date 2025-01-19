package router

import (
	"github.com/gofiber/fiber/v2"
	subQuestionChoiceHandler "github.com/gunktp20/digital-hubx-be/internal/modules/subQuestionChoice/subQuestionChoiceHandler"
)

func SetSubQuestionChoiceRoutes(api fiber.Router, subQuestionChoiceHttpHandler subQuestionChoiceHandler.SubQuestionChoiceHttpHandlerService) {
	routes := api.Group("/sub-question-choice")

	routes.Post("/", subQuestionChoiceHttpHandler.CreateSubQuestionChoice)

}
