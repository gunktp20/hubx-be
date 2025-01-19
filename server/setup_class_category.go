package server

import (
	"github.com/gofiber/fiber/v2"
	classCategoryHandler "github.com/gunktp20/digital-hubx-be/internal/modules/classCategory/classCategoryHandler"
	classCategoryRouter "github.com/gunktp20/digital-hubx-be/internal/modules/classCategory/classCategoryRouter"
	classCategoryUsecase "github.com/gunktp20/digital-hubx-be/internal/modules/classCategory/classCategoryUsecase"

	"github.com/gunktp20/digital-hubx-be/pkg/config"
)

func (s *fiberServer) initializeClassCategoryHttpHandler(api fiber.Router, conf *config.Config) {
	// ? Initialize all layers

	classCategoryUsecase := classCategoryUsecase.NewClassCategoryUsecase(s.container.Repositories.ClassCategory)
	classCategoryHttpHandler := classCategoryHandler.NewClassCategoryHttpHandler(classCategoryUsecase)

	// Routers
	classCategoryRouter.SetClassCategoryRoutes(api, classCategoryHttpHandler)
}
