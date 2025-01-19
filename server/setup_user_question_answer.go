package server

import (
	userQuestionAnswerHandler "github.com/gunktp20/digital-hubx-be/internal/modules/userQuestionAnswer/userQuestionAnswerHandler"
	userQuestionAnswerRouter "github.com/gunktp20/digital-hubx-be/internal/modules/userQuestionAnswer/userQuestionAnswerRouter"
	userQuestionAnswerUsecase "github.com/gunktp20/digital-hubx-be/internal/modules/userQuestionAnswer/userQuestionAnswerUsecase"
	"github.com/gunktp20/digital-hubx-be/pkg/config"

	"github.com/gofiber/fiber/v2"
)

func (s *fiberServer) initializeUserQuestionAnswerHttpHandler(api fiber.Router, conf *config.Config) {
	// ? Initialize all layers

	userQuestionAnswerUsecase := userQuestionAnswerUsecase.NewUserQuestionAnswerUsecase(
		s.container.Repositories.UserQuestionAnswerRepo,
		s.container.Repositories.QuestionRepo,
		s.container.Repositories.ChoiceRepo,
		s.container.Repositories.UserSubQuestionAnswerRepo,
		s.db.GetDb(),
	)
	userQuestionAnswerHttpHandler := userQuestionAnswerHandler.NewUserQuestionAnswerHttpHandler(userQuestionAnswerUsecase)

	// Routers
	userQuestionAnswerRouter.SetUserQuestionAnswerRoutes(api, userQuestionAnswerHttpHandler)
}
