package posts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/okira-e/go-as-your-backend/app/modules/users"
)

func SetupRoutes(api fiber.Router, handler *Handler, usersService *users.Service) {
	api = api.Group("/posts")

	api.Get("/", handler.FindAll)
	api.Get("/published", handler.GetPublished)
	api.Get("/count", handler.GetCount)
	api.Post(
		"/",
		users.AuthMiddleware(usersService),
		handler.Create,
	)
}
