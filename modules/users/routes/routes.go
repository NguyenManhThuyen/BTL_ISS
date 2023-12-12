package routes

import (
	usersController "app/modules/users/controller"

	"github.com/gofiber/fiber/v2"
)

func InitUsersRoutes(app *fiber.App) {
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
	api.Post("/signup", usersController.SignUpUser)
	api.Post("/signin", usersController.SignInUser)
	api.Post("/logout", usersController.LogoutUser)
}
