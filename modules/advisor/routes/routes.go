package routes

import (
	"app/modules/advisor/controller"

	"github.com/gofiber/fiber/v2"
)

func InitAdvisorRoutes(app *fiber.App) {
	advisor := app.Group("/advisor")

	getList := advisor.Group("")
	getList.Get("/", controller.GetAdvisor)
	getList.Get("/:uuid", controller.GetAdvisorByUUID)
	getList.Get("/code/:code", controller.GetAdvisorByMSCB)

	getList.Post("/create-test", controller.CreateTestAdvisors)

	getList.Post("/", controller.CreateAdvisor)
	getList.Put("/", controller.UpdateAdvisor)
	getList.Delete("/:uuid", controller.DeleteAdvisor)
	getList.Put("/restore/:uuid", controller.RestoreAdvisor)
}
