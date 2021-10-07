package api

import (
	"github.com/txfs19260817/url-shortener/database"
	"github.com/txfs19260817/url-shortener/service"
	"github.com/txfs19260817/url-shortener/utils"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/txfs19260817/url-shortener/model"
)

type request struct {
	URL      string    `json:"url" validate:"required,url"`
	Custom   string    `json:"custom,omitempty" validate:"omitempty,alphanum,max=16"`
	ExpireAt time.Time `json:"expire_at,omitempty" validate:"omitempty,gte"`
}

type response struct {
	Msg      string    `json:"msg"`
	ShortUrl string    `json:"short_url"`
	ExpireAt time.Time `json:"expire_at"`
}

// CreateUrl shorten the incoming URL
func CreateUrl(c *fiber.Ctx) error {
	// check for the incoming request body
	body := new(request)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	// add default values
	now := time.Now()
	if body.ExpireAt.IsZero() {
		body.ExpireAt = now.Add(8760 * time.Hour)
	}

	// validate request body
	if errors := utils.Validate(body); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// detect potential infinite loop
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error when getting the port we listening",
		})
	}
	if !utils.CheckNoLoopRisk(body.URL, os.Getenv("HOST"), port) {
		return c.Status(fiber.StatusLoopDetected).JSON(fiber.Map{
			"error": "haha... nice try",
		})
	}

	// return error if user's custom key has already used
	key := body.Custom
	if len(key) != 0 && database.DB.KeyExists(key) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "key exists: " + key})
	}
	// assign a hash key if user didn't give a custom one
	if len(key) == 0 {
		key = <-service.HashProvider.Queue
	}
	// do not assign an already-in-use key
	for database.DB.KeyExists(key) {
		key = <-service.HashProvider.Queue
	}

	// add to DB
	if err := database.DB.CreateUrl(&model.Url{Key: key, OriginalURL: body.URL, CreatedAt: now, ExpireAt: body.ExpireAt}); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	shortUrl := os.Getenv("HOST")
	if port != 80 && port != 443 {
		shortUrl += ":" + os.Getenv("PORT")
	}
	shortUrl += "/" + key
	return c.Status(fiber.StatusOK).JSON(&response{
		Msg:      "OK",
		ShortUrl: shortUrl,
		ExpireAt: body.ExpireAt,
	})
}
