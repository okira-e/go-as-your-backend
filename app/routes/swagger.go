package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"os"
)

// swaggerRouter sets up the swagger routes
func swaggerRouter(app *fiber.App) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app.Get("/docs/swagger.json", func(ctx *fiber.Ctx) error {
		// Only allow swagger.json to be retrieved from localhost
		if ctx.IP() != "127.0.0.1" && ctx.IP() != "::1" {
			return ctx.Status(fiber.StatusForbidden).SendString("Forbidden")
		}

		return ctx.Status(fiber.StatusOK).SendFile("./docs/swagger.json")
	})

	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "http://localhost:" + port + "/docs/swagger.json",
	}))
}
