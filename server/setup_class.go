package server

import (
	"github.com/gofiber/fiber/v2"
	classHandler "github.com/gunktp20/digital-hubx-be/internal/modules/class/classHandler"
	classRouter "github.com/gunktp20/digital-hubx-be/internal/modules/class/classRouter"
	classUsecase "github.com/gunktp20/digital-hubx-be/internal/modules/class/classUsecase"

	"github.com/gunktp20/digital-hubx-be/pkg/config"
)

func (s *fiberServer) initializeClassHttpHandler(api fiber.Router, conf *config.Config) {
	// ? Initialize all layers
	classUsecase := classUsecase.NewClassUsecase(
		s.container.Repositories.ClassRepo,
		s.container.Repositories.ClassCategory,
		s.gcs,
	)
	classHttpHandler := classHandler.NewClassHttpHandler(classUsecase)

	// Routers
	classRouter.SetClassRoutes(api, classHttpHandler)
}
