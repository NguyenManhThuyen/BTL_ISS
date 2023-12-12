package routes

import (
	"app/modules/council/controller"

	"github.com/gofiber/fiber/v2"
)

func InitCouncilRoutes(app *fiber.App) {
	// team := app.Group("/team", )
	council := app.Group("/council")

	getList := council.Group("")
	getList.Get("/", controller.GetCouncils)
	getList.Get("/code/:code", controller.GetCouncilByMSCB)
	getList.Get("/:uuid", controller.GetCouncilByUUID)

	getList.Post("/create-test", controller.CreateTestCouncil)
	getList.Post("/", controller.CreateCouncil)

	getList.Put("/", controller.UpdateCouncil)
	getList.Delete("/:uuid", controller.DeleteCouncil)
	getList.Put("/restore/:uuid", controller.RestoreCouncil)
}
