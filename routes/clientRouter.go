package routes

import (
	"github.com/gofiber/fiber/v2"
	controllers "github.com/vendz/custom-0auth/controllers"
	"github.com/vendz/custom-0auth/middleware"
)

func ClientRoutes(clientRoutes fiber.Router, h *controllers.Database) {
	clientGroup := clientRoutes.Group("/client")
	clientGroup.Use(func(c *fiber.Ctx) error {
		return middleware.UserIdInterceptor(c, h.MongoClient)
	})
	clientGroup.Post("/createClient", h.CreateClient)
}
