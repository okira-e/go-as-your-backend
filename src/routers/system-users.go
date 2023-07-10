package routers

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func SystemUsersRouter(app *fiber.App, db *sql.DB) {
	routerPath := "/api/system-users"

	// @Summary Get all system users
	// @Description Get all system users
	// @Tags Example
	// @Accept */*
	// @Produce string
	// @Success 200 {string} string "Hello, World!"
	// @Router /api/example/get [get]
	app.Get(routerPath, func(c *fiber.Ctx) error {
		fmt.Println(db.Ping())

		return c.SendString("Hello, World!")
	})

	// @Summary An example of a POST request
	// @Description An example of a POST request
	// @Tags Example
	// @Accept */*
	// @Produce string
	// @Success 200 {string} string "Hello, World 👋!"
	// @Router /api/example/post [post]
	app.Post(routerPath, func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})
}
