package router

import (
	"github.com/gofiber/fiber/v2"
	questionHandler "github.com/gunktp20/digital-hubx-be/internal/modules/question/questionHandler"
)

func SetQuestionRoutes(api fiber.Router, questionHttpHandler questionHandler.QuestionHttpHandlerService) {
	questionRoute := api.Group("/question")
	questionRoute.Get("/:class_id/class", questionHttpHandler.GetQuestionsByClassID)

	// ? Admin Routes Group
	adminRoute := api.Group("/admin/question")
	adminRoute.Post("/", questionHttpHandler.CreateQuestion)

}
