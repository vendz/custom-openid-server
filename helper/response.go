package helper

import "github.com/gofiber/fiber/v2"

func HandleError(c *fiber.Ctx, err error, status int) error {
	return c.Status(status).JSON(fiber.Map{"status": "fail", "error": err.Error()})
}

func FormatResponse(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(fiber.Map{"status": message, "data": data})
}
