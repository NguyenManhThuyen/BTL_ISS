package routes

import (
	"app/modules/facultyOffice/controller"

	"github.com/gofiber/fiber/v2"
)

func InitFacultyOfficeRoutes(app *fiber.App) {
	// team := app.Group("/team")
	facultyOffice := app.Group("/facultyoffice")

	getList := facultyOffice.Group("")
	getList.Get("/", controller.GetFacultyOffices)
	getList.Get("/code/:code", controller.GetFacultyOfficeByMSCB)
	getList.Get("/:uuid", controller.GetFacultyOfficeByUUID)

	getList.Post("/create-test", controller.CreateTestFacultyOffices)
	getList.Post("/", controller.CreateFacultyOffice)

	getList.Put("/", controller.UpdateFacultyOffice)
	getList.Delete("/:uuid", controller.DeleteFacultyOfficeUUID)
	getList.Put("/restore/:uuid", controller.RestoreFacultyOfficeUUID)
}
