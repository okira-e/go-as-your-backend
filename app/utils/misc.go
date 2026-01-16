package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/okira-e/go-as-your-backend/app/models"
)

// GetUserFromContext extracts the user from the fiber context
func GetUserFromContext(ctx *fiber.Ctx) (models.JwtUser, error) {
	user := ctx.Locals("user")
	if user == nil {
		return models.JwtUser{}, fiber.NewError(401, "User not found in context")
	}

	sessionUser, ok := user.(models.JwtUser)
	if !ok {
		return models.JwtUser{}, fiber.NewError(500, "Invalid user type in context")
	}

	return sessionUser, nil
}
