package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/txfs19260817/url-shortener/service"
)

// SetupRoutes set up a router
func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(<-service.HashProvider.Queue)
	})
	app.Post("/s", CreateUrl)
}