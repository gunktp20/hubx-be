package server

import (
	"github.com/gofiber/fiber/v2"

	subQuestionHandler "github.com/gunktp20/digital-hubx-be/internal/modules/subQuestion/subQuestionHandler"
	subQuestionRouter "github.com/gunktp20/digital-hubx-be/internal/modules/subQuestion/subQuestionRouter"
	subQuestionUsecase "github.com/gunktp20/digital-hubx-be/internal/modules/subQuestion/subQuestionUsecase"
	"github.com/gunktp20/digital-hubx-be/pkg/config"
)

func (s *fiberServer) initializeSubQuestionHttpHandler(api fiber.Router, conf *config.Config) {
	// ? Initialize all layers

	subQuestionUsecase := subQuestionUsecase.NewSubQuestionUsecase(
		s.container.Repositories.SubQuestionRepo,
		s.container.Repositories.ClassRepo,
		s.container.Repositories.ChoiceRepo)

	subQuestionHttpHandler := subQuestionHandler.NewSubQuestionHttpHandler(subQuestionUsecase)

	// Routers
	subQuestionRouter.SetSubQuestionRoutes(api, subQuestionHttpHandler)
}
