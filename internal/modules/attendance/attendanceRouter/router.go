package router

import (
	"github.com/gofiber/fiber/v2"
	attendanceHandler "github.com/gunktp20/digital-hubx-be/internal/modules/attendance/attendanceHandler"
	"github.com/gunktp20/digital-hubx-be/pkg/middleware"
)

func SetAttendanceRoute(api fiber.Router, attendanceHttpHandler attendanceHandler.AttendanceHttpHandlerService) {
	_ = api.Group("/attendance")

	// ? Admin Routes
	adminRoute := api.Group("/admin/attendance", middleware.Ident, middleware.PermissionCheck)

	// Single attendance creation
	adminRoute.Post("/", attendanceHttpHandler.CreateAttendance)
	// Multiple attendances creation
	adminRoute.Post("/batch", attendanceHttpHandler.CreateAttendances)
}
