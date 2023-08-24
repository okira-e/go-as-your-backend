package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up all the routes for the application
func SetupRoutes(app *fiber.App, db *sql.DB) {
	SwaggerRouter(app)

	PingRouter(app, db)
}
