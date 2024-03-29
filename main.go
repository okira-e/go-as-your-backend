package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"os"

	"github.com/okira-e/go-as-your-backend/app/datasource"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/okira-e/go-as-your-backend/app/routes"
	"github.com/okira-e/go-as-your-backend/app/utils"
)

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name UNLICENSE
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {
	// Setup flags
	migrate := false
	logGormTransactions := false

	for _, arg := range os.Args {
		if arg == "migrate" {
			migrate = true
		}
		if arg == "log-sql" {
			logGormTransactions = true
		}
	}

	// Load environment variables
	err := utils.LoadEnvVariables()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to database
	gormDB, err := datasource.Connect() // Doesn't report error if database is not available. We do it later.
	if err != nil {
		log.Fatal(err)
	}
	defer datasource.DisconnectOrPanic(gormDB)

	db, err := gormDB.DB()
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to database")

	// Migrate database if --migrate flag is set
	if migrate {
		err = datasource.Migrate(gormDB, logGormTransactions)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Successfully Migrated database")
		return
	}

	// Setup fiber
	app := fiber.New()

	app.Use(cors.New(
		cors.Config{
			AllowOrigins: "*",
			AllowHeaders: "Origin, Content-Type, Accept",
		},
	))

	app.Use(recover.New())

	fmt.Println("Previously generated OpenAPI docs are available at: " + os.Getenv("HOST") + ":" + os.Getenv("PORT") + "/swagger")
	fmt.Println("Running in: " + os.Getenv("APP_ENV"))

	if os.Getenv("APP_ENV") != "PROD" {
		fmt.Println("Using logger")
		app.Use(logger.New(logger.Config{
			Format: "[${ip}]:${port} ${status} - ${method} ${path} ${latency}\n",
		}))

		if logGormTransactions {
			fmt.Println("Using Gorm logger")
			gormDB.Logger = gormDB.Logger.LogMode(4)
		}
	}

	// Get the sql.DB object to pass to the routes
	sqlDB := gormDB
	if err != nil {
		log.Fatal(err)
	}

	routes.SetupRoutes(app, sqlDB)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("environment variable PORT is not set")
	}

	host := os.Getenv("HOST")
	if host == "" {
		log.Fatal("environment variable HOST is not set")
	}

	// Start server
	err = app.Listen(host + ":" + port)
	if err != nil {
		log.Fatal(err)
	}
}
