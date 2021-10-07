package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/txfs19260817/url-shortener/model"
	"time"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short,omitempty"`
	Expiry      time.Duration `json:"expiry,omitempty"`
}

type response struct {
	model.Url
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

// CreateUrl shorten the incoming URL
func CreateUrl(c *fiber.Ctx) error {
	// check for the incoming request body
	body := new(request)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}
	return nil
}