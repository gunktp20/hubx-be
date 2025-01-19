package server

import (
	"github.com/gofiber/fiber/v2"
	classRegistrationHandler "github.com/gunktp20/digital-hubx-be/internal/modules/classRegistration/classRegistrationHandler"
	classRegistrationRouter "github.com/gunktp20/digital-hubx-be/internal/modules/classRegistration/classRegistrationRouter"
	classRegistrationUsecase "github.com/gunktp20/digital-hubx-be/internal/modules/classRegistration/classRegistrationUsecase"

	"github.com/gunktp20/digital-hubx-be/pkg/config"
)

func (s *fiberServer) initializeClassRegistrationHttpHandler(api fiber.Router, conf *config.Config) {
	// ? Initialize all layers
	classRegistrationUsecase := classRegistrationUsecase.NewClassRegistrationUsecase(
		s.container.Repositories.ClassRegistrationRepo,
		s.container.Repositories.ClassSessionRepo,
		s.container.Repositories.ClassRepo,
	)
	classRegistrationHttpHandler := classRegistrationHandler.NewClassRegistrationHttpHandler(classRegistrationUsecase)

	// Routers
	classRegistrationRouter.SetClassRegistrationRoutes(api, classRegistrationHttpHandler)
}
