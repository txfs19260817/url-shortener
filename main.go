package main

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"
	"github.com/txfs19260817/url-shortener/api"
	"github.com/txfs19260817/url-shortener/config"
	"github.com/txfs19260817/url-shortener/database"
	"github.com/txfs19260817/url-shortener/service"
)

var Log = logrus.New()

func main() {
	var err error
	// log
	Log.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"}) // Log as JSON instead of the default ASCII formatter.
	Log.SetOutput(os.Stdout)                                                        // Output to stdout instead of the default stderr
	Log.SetLevel(logrus.InfoLevel)                                                  // Only Log the info severity or above.
	Log.SetReportCaller(true)                                                       // Trace callers
	w := Log.Writer()
	defer w.Close()
	// config
	c := &config.Config{}
	if err := c.ReadConfig("./config/config.yaml"); err != nil {
		Log.WithField("error", err).Fatal("error in reading config file")
	}
	// config - database
	database.DB, err = database.NewMongoDB(c.Mongodb)
	if err != nil {
		Log.WithField("error", err).Fatal("error in connecting to database")
	}
	// config - service - hashGenerator
	service.HashProvider = service.NewHashGenerator(c.Service.HashLen, c.Service.PoolSize, time.Duration(c.Service.Duration)*time.Second)
	go func() {
		err := service.HashProvider.NanoIDProvider()
		if err != nil {
			Log.WithField("error", err).Fatal("error in HashGenerator: NanoID Provider")
		}
	}()
	defer service.HashProvider.CloseProvider()

	// fiber
	app := fiber.New()
	app.Use(logger.New(logger.Config{Output: w}))
	api.SetupRoutes(app)
	Log.Fatal(app.Listen(":3001"))
}
