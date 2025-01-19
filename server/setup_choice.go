package server

import (
	"github.com/gofiber/fiber/v2"
	choiceHandler "github.com/gunktp20/digital-hubx-be/internal/modules/choice/choiceHandler"
	choiceRouter "github.com/gunktp20/digital-hubx-be/internal/modules/choice/choiceRouter"
	choiceUsecase "github.com/gunktp20/digital-hubx-be/internal/modules/choice/choiceUsecase"

	"github.com/gunktp20/digital-hubx-be/pkg/config"
)

func (s *fiberServer) initializeChoiceHttpHandler(api fiber.Router, conf *config.Config) {
	// ? Initialize all layers

	choiceUsecase := choiceUsecase.NewChoiceUsecase(
		s.container.Repositories.ChoiceRepo,
		s.container.Repositories.QuestionRepo,
	)
	choiceHttpHandler := choiceHandler.NewChoiceHttpHandler(choiceUsecase)

	// Routers
	choiceRouter.SetChoiceRoutes(api, choiceHttpHandler)
}
