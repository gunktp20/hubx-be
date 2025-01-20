package router

import (
	"github.com/gofiber/fiber/v2"
	userQuestionAnswerHandler "github.com/gunktp20/digital-hubx-be/internal/modules/userQuestionAnswer/userQuestionAnswerHandler"
)

func SetUserQuestionAnswerRoutes(api fiber.Router, userQuestionAnswerHttpHandler userQuestionAnswerHandler.UserQuestionAnswerHttpHandlerService) {
	userQuestionAnswerRoute := api.Group("/user-question-answer")
	userQuestionAnswerRoute.Get("/:class_id/class", userQuestionAnswerHttpHandler.GetUserQuestionAnswersWithClassId)

	// ? Admin Routes Group
	adminRoute := api.Group("/admin/class-session")
	adminRoute.Post("/:class_id/class", userQuestionAnswerHttpHandler.CreateMultipleUserQuestionAnswers)

}
