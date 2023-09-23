package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/vendz/custom-0auth/database"
	"github.com/vendz/custom-0auth/helper"
	"github.com/vendz/custom-0auth/routes"
)

func main() {
	app := fiber.New()

	helper.LoadEnv()
	handler := database.NewDatabase()

	app.Use(cors.New())
	app.Use(logger.New())

	apiGroup := app.Group("/api/v1")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "API is Up and Running... 🚀", "status": "success"})
	})
	routes.UserRoutes(apiGroup, &handler)
	routes.ClientRoutes(app, &handler)

	err := app.Listen(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}
}
