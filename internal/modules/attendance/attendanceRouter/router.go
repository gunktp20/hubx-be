package router

import (
	"github.com/gofiber/fiber/v2"
	attendanceHandler "github.com/gunktp20/digital-hubx-be/internal/modules/attendance/attendanceHandler"
)

func SetAttendanceRoute(api fiber.Router, attendanceHttpHandler attendanceHandler.AttendanceHttpHandlerService) {
	_ = api.Group("/attendance")

	// ? Admin Routes Group
	adminRoute := api.Group("/admin/attendance")
	adminRoute.Post("/", attendanceHttpHandler.CreateAttendance)

}
