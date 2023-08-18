package middleware

import (
	"context"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/vendz/custom-0auth/models"
	"github.com/vendz/custom-0auth/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func VerifyTokenAndDb(c *fiber.Ctx, mongoClient *mongo.Client, redisClient *redis.Client) error {

	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "error": "Unauthorized"})
	}

	fields := strings.Fields(authHeader)
	if len(fields) < 2 || fields[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "error": "Unauthorized"})
	}

	token := fields[1]
	email, err := utils.ValidateToken(token, os.Getenv("JWT_SECRET"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "error": "Unauthorized"})
	}

	filter := bson.M{"email": email}
	userCollection := mongoClient.Database("blabber").Collection("users")
	var user models.User
	err = userCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "error": "Unauthorized"})
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	rdbStr := "token:" + email
	val, err := redisClient.Get(context.Background(), rdbStr).Result()
	if err != nil {
		if err != redis.Nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token"})
	}

	if val == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token"})
	}

	if val != token {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token"})
	}

	c.Locals("email", email)
	c.Locals("user", user.Id.Hex())
	c.Locals("token", token)
	return c.Next()
}
