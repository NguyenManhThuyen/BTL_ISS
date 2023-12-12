package controller

import (
	"app/config"

	"github.com/gofiber/fiber/v2"
)



// Check token ================================================================================
func CheckToken(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	response.Status = true

	// Return response
	return c.JSON(response)
}
