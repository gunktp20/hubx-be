package router

import (
	"github.com/gofiber/fiber/v2"
	userQuestionAnswerHandler "github.com/gunktp20/digital-hubx-be/internal/modules/userQuestionAnswer/userQuestionAnswerHandler"
	"github.com/gunktp20/digital-hubx-be/pkg/middleware"
)

func SetUserQuestionAnswerRoutes(api fiber.Router, userQuestionAnswerHttpHandler userQuestionAnswerHandler.UserQuestionAnswerHttpHandlerService) {
	userQuestionAnswerRoute := api.Group("/user-question-answer")
	userQuestionAnswerRoute.Get("/:class_id/class", userQuestionAnswerHttpHandler.GetUserQuestionAnswersWithClassId)
	userQuestionAnswerRoute.Post("/:class_id/class", userQuestionAnswerHttpHandler.CreateMultipleUserQuestionAnswers)

	// ? Admin Routes
	_ = api.Group("/admin/class-session", middleware.Ident, middleware.PermissionCheck)

}
