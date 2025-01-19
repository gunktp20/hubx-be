package router

import (
	"github.com/gofiber/fiber/v2"
	classHandler "github.com/gunktp20/digital-hubx-be/internal/modules/class/classHandler"
	"github.com/gunktp20/digital-hubx-be/pkg/middleware"
)

func SetClassRoutes(api fiber.Router, classHttpHandler classHandler.ClassHttpHandlerService) {
	routes := api.Group("/class", middleware.DecodeJwtPayload)

	routes.Get("/", classHttpHandler.GetAllClasses)
	routes.Post("/", classHttpHandler.CreateClass)
	routes.Get("/:class_id", classHttpHandler.GetClassById)

}
