package router

import (
	"github.com/gofiber/fiber/v2"
	subQuestionHandler "github.com/gunktp20/digital-hubx-be/internal/modules/subQuestion/subQuestionHandler"
)

func SetSubQuestionRoutes(api fiber.Router, subQuestionHttpHandler subQuestionHandler.SubQuestionHttpHandlerService) {
	routes := api.Group("/sub-question")

	routes.Post("/", subQuestionHttpHandler.CreateSubQuestion)
	routes.Get("/:question_id/question", subQuestionHttpHandler.GetSubQuestionsByQuestionID)
	routes.Get("/:choice_id/choice", subQuestionHttpHandler.GetSubQuestionsByChoiceID)

}
