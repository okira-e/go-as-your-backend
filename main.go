package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/okira-e/go-as-your-backend/app/modules/posts"
	"github.com/okira-e/go-as-your-backend/app/modules/roles"
	"github.com/okira-e/go-as-your-backend/app/modules/users"
	"github.com/okira-e/go-as-your-backend/app/utils"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()

	env := utils.RequireEnv("ENV")

	db, err := setupDatabase()
	if err != nil {
		log.Fatalf("Error connecting to the database. %s\n", err.Error())
	}

	app := fiber.New()

	// Setup CORS.
	clientOrigin := utils.RequireEnv("CLIENT_URL")

	app.Use(cors.New(cors.Config{
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
		AllowOrigins:     clientOrigin,
	}))

	// For recovering from panics
	app.Use(recover.New())

	// Setup logging
	if env == "debug" {
		app.Use(logger.New(logger.Config{
			Format: "[${ip}]:${port} ${status} - ${method} ${path} ${latency}\n",
		}))
	} else {
		env = "release"
		// Setup optional writing to persistent logs.
	}

	setupModules(app, db)

	fmt.Printf("Running in %s environment.\n", env)

	port := utils.RequireEnv("PORT")
	host := utils.RequireEnv("HOST")

	log.Fatalln(app.Listen(host + ":" + port))
}

func setupDatabase() (*gorm.DB, error) {
	host := utils.RequireEnv("DB_HOST")
	port := utils.RequireEnv("DB_PORT")
	user := utils.RequireEnv("DB_USER")
	dbname := utils.RequireEnv("DB_NAME")
	password := utils.RequireEnv("DB_PASSWORD")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func setupModules(app *fiber.App, db *gorm.DB) {
	version := utils.RequireEnv("API_VERSION")

	api := app.Group("/api")

	// Health endpoint
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	versionedApi := api.Group("/" + version)

	usersRepo := users.NewRepository(db)
	usersService := users.NewService(usersRepo)
	usersHandler := users.NewHandler(usersService)
	users.SetupRoutes(versionedApi, usersHandler, usersService)

	rolesRepo := roles.NewRepository(db)
	rolesService := roles.NewService(rolesRepo)
	rolesHandler := roles.NewHandler(rolesService)
	roles.SetupRoutes(versionedApi, rolesHandler)

	postsRepo := posts.NewRepository(db)
	postsService := posts.NewService(postsRepo)
	postsHandler := posts.NewHandler(postsService)
	posts.SetupRoutes(versionedApi, postsHandler, usersService)
}
