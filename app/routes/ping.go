package routes

import (
	"github.com/gofiber/fiber/v2"
)

func pingRouter(app *fiber.App) {
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
