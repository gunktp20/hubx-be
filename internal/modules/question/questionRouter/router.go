package router

import (
	"github.com/gofiber/fiber/v2"
	questionHandler "github.com/gunktp20/digital-hubx-be/internal/modules/question/questionHandler"
	"github.com/gunktp20/digital-hubx-be/pkg/middleware"
)

func SetQuestionRoutes(api fiber.Router, questionHttpHandler questionHandler.QuestionHttpHandlerService) {
	questionRoute := api.Group("/question")
	questionRoute.Get("/:class_id/class", questionHttpHandler.GetQuestionsByClassID)

	// ? Admin Routes
	adminRoute := api.Group("/admin/question", middleware.PermissionCheck)
	adminRoute.Post("/", questionHttpHandler.CreateQuestion)

}
