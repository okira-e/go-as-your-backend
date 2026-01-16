package users

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(api fiber.Router, handler *Handler, usersService *Service) {
	api = api.Group("/users")

	api.Post("/login", handler.Login)
	api.Post("/register", handler.Register)
	api.Post("/refresh", handler.RefreshToken)
	api.Post("/logout", handler.Logout)
	api.Get("/validate-token", handler.ValidateToken)

	api.Get("/me", AuthMiddleware(usersService), handler.Me)
	api.Get("/", AuthMiddleware(usersService), handler.FindAll)
	// @TODO: Know how to secure this as it now doxes user info w/out auth
	api.Get("/contact-info/:id", handler.GetContactInfo)
	api.Get("/count", AuthMiddleware(usersService), handler.GetCount)

	// users.Post("/", handler.CreateUser)
}
