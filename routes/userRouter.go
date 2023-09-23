package routes

import (
	"github.com/gofiber/fiber/v2"
	controllers "github.com/vendz/custom-0auth/controllers"
	"github.com/vendz/custom-0auth/middleware"
)

func UserRoutes(incomingRoutes fiber.Router, h *controllers.Database) {
	authGroup := incomingRoutes.Group("/auth")
	authGroup.Use(middleware.RedirectInterceptor, middleware.ClientIdInterceptor)
	authGroup.Post("/createUser", h.CreateUser)
	authGroup.Post("/login", h.LoginUser)

	userGroup := incomingRoutes.Group("/user")
	userGroup.Use(func(c *fiber.Ctx) error {
		return middleware.VerifyTokenAndDb(c, h.MongoClient, h.RedisClient)
	})
	userGroup.Post("/logout", h.LogoutUser)
	userGroup.Get("/me", h.GetMe)
	userGroup.Get("/sso", middleware.RedirectInterceptor, middleware.ClientIdInterceptor, h.SingleSignon)
}
