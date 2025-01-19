package server

import (
	"github.com/gofiber/fiber/v2"

	subQuestionChoiceHandler "github.com/gunktp20/digital-hubx-be/internal/modules/subQuestionChoice/subQuestionChoiceHandler"
	subQuestionChoiceRouter "github.com/gunktp20/digital-hubx-be/internal/modules/subQuestionChoice/subQuestionChoiceRouter"
	subQuestionChoiceUsecase "github.com/gunktp20/digital-hubx-be/internal/modules/subQuestionChoice/subQuestionChoiceUsecase"
	"github.com/gunktp20/digital-hubx-be/pkg/config"
)

func (s *fiberServer) initializeSubQuestionChoiceHttpHandler(api fiber.Router, conf *config.Config) {
	// ? Initialize all layers

	subQuestionChoiceUsecase := subQuestionChoiceUsecase.NewSubQuestionChoiceUsecase(
		s.container.Repositories.SubQuestionChoiceRepo,
		s.container.Repositories.SubQuestionRepo,
	)

	subQuestionChoiceHttpHandler := subQuestionChoiceHandler.NewSubQuestionChoiceHttpHandler(subQuestionChoiceUsecase)

	// Routers
	subQuestionChoiceRouter.SetSubQuestionChoiceRoutes(api, subQuestionChoiceHttpHandler)
}
