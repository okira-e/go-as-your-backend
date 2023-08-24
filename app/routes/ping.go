package routes

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
)

func PingRouter(app *fiber.App, _ *sql.DB) {
	routerPath := "/api/ping"

	app.Get(routerPath, pingServer)
}

// @Summary Ping the server
// @Description Ping the server
// @Tags ping
// @Accept */*
// @Produce plain
// @Success 200 "OK"
// @Router /api/ping [get]
func pingServer(c *fiber.Ctx) error {
	return c.SendString("pong")
}
