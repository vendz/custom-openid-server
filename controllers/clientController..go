package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vendz/custom-0auth/helper"
)

func (databaseClient Database) CreateClient(c *fiber.Ctx) error {
	return helper.FormatResponse(c, fiber.StatusOK, "success", fiber.Map{"message": "client created!", "user": c.Locals("user")})
}
