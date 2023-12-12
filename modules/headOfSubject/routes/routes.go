package routes

import (
	"app/modules/headOfSubject/controller"

	"github.com/gofiber/fiber/v2"
)

func InitHeadOfSubjectRoutes(app *fiber.App) {
	headOfSubject := app.Group("/headofsubject")

	getList := headOfSubject.Group("")
	getList.Get("/", controller.GetHeadOfSubject)
	getList.Get("/:uuid", controller.GetHeadOfSubjectByUUID)
	getList.Get("/code/:code", controller.GetHeadOfSubjectByMSCB)

	getList.Post("/create-test", controller.CreateTestHeadOfSubjects)

	getList.Post("/", controller.CreateHeadOfSubject)
	getList.Put("/", controller.UpdateHeadOfSubject)
	getList.Delete("/:uuid", controller.DeleteHeadOfSubjectByUUID)
	getList.Put("/restore/:uuid", controller.RestoreHeadOfSubjectByUUID)
}
