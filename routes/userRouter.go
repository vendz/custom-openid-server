package routes

import (
	"github.com/gofiber/fiber/v2"
	controllers "github.com/vendz/custom-0auth/controllers"
	"github.com/vendz/custom-0auth/middleware"
)

func UserRoutes(incomingRoutes *fiber.App, h *controllers.Database) {
	incomingRoutes.Post("/api/v1/createUser", h.CreateUser)
	incomingRoutes.Post("/api/v1/login", h.LoginUser)

	userGroup := incomingRoutes.Group("/api/v1/user", func(c *fiber.Ctx) error {
		err := middleware.VerifyTokenAndDb(c, h.MongoClient, h.RedisClient)
		if err != nil {
			return err
		}
		return nil
	})
	userGroup.Post("/logout", h.LogoutUser)
	userGroup.Get("/me", h.GetMe)
	userGroup.Get("/authenticate", h.AuthenticateUser)
}
