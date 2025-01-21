package router

import (
	"github.com/gofiber/fiber/v2"
	classRegistrationHandler "github.com/gunktp20/digital-hubx-be/internal/modules/classRegistration/classRegistrationHandler"
	"github.com/gunktp20/digital-hubx-be/pkg/middleware"
)

func SetClassRegistrationRoutes(api fiber.Router, classRegistrationHttpHandler classRegistrationHandler.ClassRegistrationHttpHandlerService) {
	classRegistrationRoute := api.Group("/class-registration")

	classRegistrationRoute.Get("/", classRegistrationHttpHandler.GetUserRegistrations)
	classRegistrationRoute.Post("/", classRegistrationHttpHandler.CreateClassRegistration)
	classRegistrationRoute.Delete("/:class_session_id/cancel", classRegistrationHttpHandler.CancelClassRegistration)

	// ? Admin Routes Group
	adminRoute := api.Group("/admin/class-registration", middleware.Ident, middleware.PermissionCheck)
	adminRoute.Post("/reset-cancel-quota", classRegistrationHttpHandler.ResetCancelledQuota)
	adminRoute.Delete("/:class_session_id/:email", classRegistrationHttpHandler.DeleteUserClassRegistrationBySession)

}
