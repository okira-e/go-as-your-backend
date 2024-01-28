package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var db *gorm.DB

// SetupRoutes sets up all the routes for the application
func SetupRoutes(app *fiber.App, datasource *gorm.DB) {
	db = datasource

	swaggerRouter(app)
	pingRouter(app)
	systemUsersRouter(app)
}
