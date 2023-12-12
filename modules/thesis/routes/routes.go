package routes

import (
	"app/modules/thesis/controller"

	"github.com/gofiber/fiber/v2"
)

func InitThesisRoutes(app *fiber.App) {
	// Define a group for /thesis route with middleware
	thesis := app.Group("/thesis")

	// Define your thesis API routes
	thesis.Get("/", controller.ListTheses)

	thesis.Get("/get-by-createby/{createBy}", controller.GetThesesByCreateBy)
	thesis.Get("/:uuid", controller.GetThesis)

	thesis.Post("/", controller.CreateThesis)
	thesis.Post("/create-test", controller.CreateTestTheses)
	thesis.Post("/status-thesis", controller.PostStatusThesis)
	thesis.Put("/", controller.UpdateThesis)
	thesis.Put("/approval",controller.UpdateThesisApprovalStatus)
	thesis.Delete("/:uuid", controller.DeleteThesis)

	// Additional routes for adding and removing students and advisors
	thesis.Post("/addstudent/:thesisUUID/:studentUUID", controller.AddStudentToThesis)
	thesis.Delete("/removestudent/:thesisUUID/:studentUUID", controller.RemoveStudentFromThesis)
	thesis.Post("/addadvisor/:thesisUUID/:advisorUUID", controller.AddAdvisorToThesis)
	thesis.Delete("/removeadvisor/:thesisUUID/:advisorUUID", controller.RemoveAdvisorFromThesis)

}
