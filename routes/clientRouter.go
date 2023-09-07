package routes

import (
	"github.com/gofiber/fiber/v2"
	controllers "github.com/vendz/custom-0auth/controllers"
	"github.com/vendz/custom-0auth/middleware"
)

func ClientRoutes(clientRoutes *fiber.App, h *controllers.Database) {
	clientGroup := clientRoutes.Group("/client", func(c *fiber.Ctx) error {
		err := middleware.UserIdInterceptor(c, h.MongoClient)
		if err != nil {
			return err
		}
		return nil
	})
	clientGroup.Post("/createClient", h.CreateClient)
}
