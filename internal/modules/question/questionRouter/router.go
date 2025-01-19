package router

import (
	"github.com/gofiber/fiber/v2"
	questionHandler "github.com/gunktp20/digital-hubx-be/internal/modules/question/questionHandler"
)

func SetQuestionRoutes(api fiber.Router, questionHttpHandler questionHandler.QuestionHttpHandlerService) {
	routes := api.Group("/question")

	routes.Post("/", questionHttpHandler.CreateQuestion)
	routes.Get("/:class_id/class", questionHttpHandler.GetQuestionsByClassID)

}
