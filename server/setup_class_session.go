package server

import (
	"github.com/gofiber/fiber/v2"
	classSessionHandler "github.com/gunktp20/digital-hubx-be/internal/modules/classSession/classSessionHandler"
	classSessionRouter "github.com/gunktp20/digital-hubx-be/internal/modules/classSession/classSessionRouter"
	classSessionUsecase "github.com/gunktp20/digital-hubx-be/internal/modules/classSession/classSessionUsecase"

	"github.com/gunktp20/digital-hubx-be/pkg/config"
)

func (s *fiberServer) initializeClassSessionHttpHandler(api fiber.Router, conf *config.Config) {
	// ? Initialize all layers
	classSessionUsecase := classSessionUsecase.NewClassSessionUsecase(
		s.container.Repositories.ClassSessionRepo,
		s.container.Repositories.ClassRepo,
		s.gcs,
	)
	classSessionHttpHandler := classSessionHandler.NewClassSessionHttpHandler(classSessionUsecase)

	// Routers
	classSessionRouter.SetClassSessionRoutes(api, classSessionHttpHandler)
}
