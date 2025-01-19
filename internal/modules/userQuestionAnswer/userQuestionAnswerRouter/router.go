package router

import (
	"github.com/gofiber/fiber/v2"
	userQuestionAnswerHandler "github.com/gunktp20/digital-hubx-be/internal/modules/userQuestionAnswer/userQuestionAnswerHandler"
)

func SetUserQuestionAnswerRoutes(api fiber.Router, userQuestionAnswerHttpHandler userQuestionAnswerHandler.UserQuestionAnswerHttpHandlerService) {
	routes := api.Group("/user-question-answer")

	// routes.Post("/bulk", userQuestionAnswerHttpHandler.CreateMultipleUserQuestionAnswers)
	routes.Post("/:class_id/class", userQuestionAnswerHttpHandler.CreateMultipleUserQuestionAnswers)
	routes.Get("/:class_id/class", userQuestionAnswerHttpHandler.GetUserQuestionAnswersWithClassId)

}
