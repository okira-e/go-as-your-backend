package roles

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(api fiber.Router, handler *Handler) {
	api = api.Group("/roles")

	api.Get("/", handler.FindAll)
	api.Get("/count", handler.GetCount)
	api.Post("/", handler.Create)
}
