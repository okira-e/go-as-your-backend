package utils

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool   `json:"success"`
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Err(ctx *fiber.Ctx, status int, message string, data any) error {
	return ctx.Status(status).JSON(Response{
		Success: false,
		Status:  status,
		Message: message,
		Data:    data,
	})
}

func Ok(ctx *fiber.Ctx, status int, message string, data any) error {
	return ctx.Status(status).JSON(Response{
		Success: true,
		Status:  status,
		Message: message,
		Data:    data,
	})
}
