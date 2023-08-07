package response

import "github.com/gofiber/fiber/v2"

func Success(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   data,
	})
}

func Error(c *fiber.Ctx, statusCode int, errorMessage string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"status":  "error",
		"message": errorMessage,
	})
}
