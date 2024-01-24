package utils

import (
	"github.com/okira-e/go-as-your-backend/app/types"
	"github.com/gofiber/fiber/v2"
)

func Err(c *fiber.Ctx, code int, message string, data any) error {
	response := types.Response{
		Success: false,
		Code:    code,
		Message: message,
		Data:    data,
	}

	return c.Status(code).JSON(response)
}

func Ok(c *fiber.Ctx, code int, message string, data any) error {
	response := types.Response{
		Success: true,
		Code:    code,
		Message: message,
		Data:    data,
	}

	return c.Status(code).JSON(response)
}

