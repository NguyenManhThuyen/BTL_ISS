package moduleauthen

import (
	authenController "app/modules/authen/controller"

	"github.com/gofiber/fiber/v2"
)

func InitAuthenRoutes(app *fiber.App) {
	/**
	*
	*	System User
	* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	*
	**/
	api := app.Group("/")

	/**
	*
	*	Authen
	* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	*
	**/
	//api.Post("/login", authenController.Login)
	api.Post("/check-token", authenController.CheckToken)
}
