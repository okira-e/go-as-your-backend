package routes

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/okira-e/go-as-your-backend/app/utils"
)

func PingRouter(app *fiber.App) {
	routerPath := "/api/ping"

	app.Get(routerPath, pingServer)

	app.Post(routerPath, exampleRoute)
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

func exampleRoute(c *fiber.Ctx) error {
	body := struct {
		Name string `json:"name"`
		Age  uint   `json:"age"`
	}{}

	err := json.Unmarshal(c.Body(), &body)
	if err != nil {
		return err
	}

	// Check if the `name` and `age` fields were provided in the request body.
	// If not, "ValidateFields" will send a response of 400 and a message with
	// the name of the fields not provided.
	if ret, pass := utils.ValidateFields(c, body, "Name", "Age"); !pass {
		return ret
	}

	return utils.Ok(c, 200, "OK", nil)
}
