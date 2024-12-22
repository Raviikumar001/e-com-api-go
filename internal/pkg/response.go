package response

import "github.com/gofiber/fiber/v2"

type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func Success(c *fiber.Ctx, data interface{}) error {
    return c.JSON(Response{
        Success: true,
        Data:    data,
    })
}

func Error(c *fiber.Ctx, statusCode int, message string) error {
    return c.Status(statusCode).JSON(Response{
        Success: false,
        Error:   message,
    })
}