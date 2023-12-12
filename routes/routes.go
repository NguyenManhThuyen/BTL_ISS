package routes

import (

	_"app/docs"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {

	// Document
	app.Get("/document/*", swagger.HandlerDefault)

}
