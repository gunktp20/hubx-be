package server

import (
	"github.com/gofiber/fiber/v2"

	questionHandler "github.com/gunktp20/digital-hubx-be/internal/modules/question/questionHandler"
	questionRouter "github.com/gunktp20/digital-hubx-be/internal/modules/question/questionRouter"
	questionUsecase "github.com/gunktp20/digital-hubx-be/internal/modules/question/questionUsecase"
	"github.com/gunktp20/digital-hubx-be/pkg/config"
)

func (s *fiberServer) initializeQuestionHttpHandler(api fiber.Router, conf *config.Config) {
	// ? Initialize all layers
	questionUsecase := questionUsecase.NewQuestionUsecase(
		s.container.Repositories.QuestionRepo,
		s.container.Repositories.ClassRepo,
	)
	questionHttpHandler := questionHandler.NewQuestionHttpHandler(questionUsecase)

	// Routers
	questionRouter.SetQuestionRoutes(api, questionHttpHandler)
}
