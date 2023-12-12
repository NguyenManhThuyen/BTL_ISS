package routes

import (
	"app/modules/student/controller"

	"github.com/gofiber/fiber/v2"
)

func InitStudentRoutes(app *fiber.App) {
	// team := app.Group("/team")
	student := app.Group("/student")

	getList := student.Group("")
	getList.Get("/", controller.GetStudent)
	getList.Get("/code/:code", controller.GetStudentByCode)
	getList.Get("/:uuid", controller.GetStudentByUUID)

	getList.Post("/create-test", controller.CreateTestStudents)
	getList.Post("/", controller.CreateStudent)

	getList.Put("/", controller.UpdateStudent)
	getList.Delete("/:uuid", controller.DeleteStudentByUUID)
	getList.Put("/restore/:uuid", controller.RestoreStudentByUUID)
}
