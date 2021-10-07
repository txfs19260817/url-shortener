package main

import (
	"fmt"
	"github.com/txfs19260817/url-shortener/utils"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/txfs19260817/url-shortener/api"
	"github.com/txfs19260817/url-shortener/config"
	"github.com/txfs19260817/url-shortener/database"
	"github.com/txfs19260817/url-shortener/service"
)

func main() {
	var err error
	// log
	w := utils.InitLogger()
	defer w.Close()
	// config
	c := &config.Config{}
	if err := c.ReadConfig("./config/config.yaml"); err != nil {
		utils.Log.WithField("error", err).Fatal("error in reading config file")
	}
	// config - database
	database.DB, err = database.NewMongoDB(c.Mongodb)
	if err != nil {
		utils.Log.WithField("error", err).Fatal("error in connecting to database")
	}
	// config - service - hashGenerator
	service.HashProvider = service.NewHashGenerator(c.Service.HashLen, c.Service.PoolSize, time.Duration(c.Service.Duration)*time.Second)
	go func() {
		if err := service.HashProvider.NanoIDProvider(); err != nil {
			utils.Log.WithField("error", err).Fatal("error in HashGenerator: NanoID Provider")
		}
	}()
	defer service.HashProvider.CloseProvider()

	// fiber
	app := fiber.New()
	app.Use(logger.New(logger.Config{Output: w}))
	api.SetupRoutes(app)
	if err := os.Setenv("HOST", c.Server.Host); err != nil {
		utils.Log.WithField("error", err).WithField("env", c.Server.Host).Fatal("error in Setenv: HOST")
	}
	if err := os.Setenv("PORT", strconv.Itoa(c.Server.Port)); err != nil {
		utils.Log.WithField("error", err).WithField("env", c.Server.Port).Fatal("error in Setenv: PORT")
	}
	utils.Log.Fatal(app.Listen(fmt.Sprintf(c.Server.Host+":%d", c.Server.Port)))
}
