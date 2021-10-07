package api

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes set up a router
func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hi")
	})
	app.Post("/s", CreateUrl)
}
