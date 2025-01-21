package router

import (
	"github.com/gofiber/fiber/v2"
	subQuestionHandler "github.com/gunktp20/digital-hubx-be/internal/modules/subQuestion/subQuestionHandler"
	"github.com/gunktp20/digital-hubx-be/pkg/middleware"
)

func SetSubQuestionRoutes(api fiber.Router, subQuestionHttpHandler subQuestionHandler.SubQuestionHttpHandlerService) {
	subQuestionRoute := api.Group("/sub-question")

	subQuestionRoute.Get("/:question_id/question", subQuestionHttpHandler.GetSubQuestionsByQuestionID)
	subQuestionRoute.Get("/:choice_id/choice", subQuestionHttpHandler.GetSubQuestionsByChoiceID)

	// ? Admin Routes Group
	adminRoute := api.Group("/admin/class-session", middleware.PermissionCheck)
	adminRoute.Post("/", subQuestionHttpHandler.CreateSubQuestion)

}
