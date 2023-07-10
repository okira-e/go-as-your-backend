package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Okira-E/go-as-your-backend/src/datasource"

	"github.com/Okira-E/go-as-your-backend/src/routers"
	"github.com/Okira-E/go-as-your-backend/src/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// -- Setup flags
	migrate := false

	for _, arg := range os.Args {
		if arg == "migrate" {
			migrate = true
		}
	}

	// -- Load environment variables
	err := utils.LoadEnvVariables()
	if err != nil {
		log.Fatal(err)
	}

	// -- Connect to database
	gormDB, err := datasource.Connect()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connected to database")
	}
	defer datasource.DisconnectOrPanic(gormDB)

	// -- Migrate database if --migrate flag is set
	if migrate {
		err = datasource.Migrate(gormDB)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Migrated database")
		return
	}

	// -- Setup fiber
	app := fiber.New()

	app.Use(recover.New())

	fmt.Println("Running in: " + os.Getenv("APP_ENV"))

	if os.Getenv("ENV") != "PROD" {
		app.Use(logger.New(logger.Config{
			Format: "[${ip}]:${port} ${status} - ${method} ${path} ${latency}\n",
		}))
	}

	// -- Setup routes
	// Get the sql.DB object to pass to the routes
	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatal(err)
	}

	routers.SetupRoutes(app, sqlDB)

	// -- Start server
	port := os.Getenv("PORT")
	err = app.Listen(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}
