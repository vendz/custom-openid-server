package routes

import (
	"github.com/gofiber/fiber/v2"
	controllers "github.com/vendz/custom-0auth/controllers"
	"github.com/vendz/custom-0auth/middleware"
)

func UserRoutes(incomingRoutes *fiber.App, h *controllers.Database) {
	authGroup := incomingRoutes.Group("/api/v1/auth", middleware.RedirectInterceptor, middleware.ClientIdInterceptor)
	authGroup.Post("/createUser", h.CreateUser)
	authGroup.Post("/login", h.LoginUser)

	userGroup := incomingRoutes.Group("/api/v1/user", func(c *fiber.Ctx) error {
		err := middleware.VerifyTokenAndDb(c, h.MongoClient, h.RedisClient)
		if err != nil {
			return err
		}
		return nil
	})
	userGroup.Post("/logout", h.LogoutUser)
	userGroup.Get("/me", h.GetMe)

	ssoGroup := incomingRoutes.Group("/api/v1", func(c *fiber.Ctx) error {
		err := middleware.VerifyTokenAndDb(c, h.MongoClient, h.RedisClient)
		if err != nil {
			return err
		}
		return nil
	})
	ssoGroup.Get("/sso", middleware.RedirectInterceptor, middleware.ClientIdInterceptor, h.SingleSignon)

	// clientGroup := incomingRoutes.Group("/api/v1/user", func(c *fiber.Ctx) error {
	// 	err := middleware.VerifyTokenAndDb(c, h.MongoClient, h.RedisClient)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	return nil
	// })
	// clientGroup.Post("/createClient", h.AuthenticateUser)
}
