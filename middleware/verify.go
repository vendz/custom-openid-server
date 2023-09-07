package middleware

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/vendz/custom-0auth/helper"
	"github.com/vendz/custom-0auth/models"
	"github.com/vendz/custom-0auth/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func VerifyTokenAndDb(c *fiber.Ctx, mongoClient *mongo.Client, redisClient *redis.Client) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return helper.HandleError(c, fmt.Errorf("no authorization header found"), fiber.StatusUnauthorized)
	}

	fields := strings.Fields(authHeader)
	if len(fields) < 2 || fields[0] != "Bearer" {
		return helper.HandleError(c, fmt.Errorf("token not fount"), fiber.StatusUnauthorized)
	}

	token := fields[1]
	email, err := utils.ValidateToken(token, os.Getenv("JWT_SECRET"))
	if err != nil {
		return helper.HandleError(c, fmt.Errorf("invalid token"), fiber.StatusUnauthorized)
	}

	filter := bson.M{"email": email}
	userCollection := mongoClient.Database("custom-auth").Collection("users")
	var user models.User
	err = userCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return helper.HandleError(c, fmt.Errorf("user not found"), fiber.StatusUnauthorized)
		}
		return helper.HandleError(c, err, fiber.StatusInternalServerError)
	}

	rdbStr := "token:" + email
	val, err := redisClient.Get(context.Background(), rdbStr).Result()
	if err != nil {
		if err != redis.Nil {
			return helper.HandleError(c, err, fiber.StatusInternalServerError)
		}
		return helper.HandleError(c, fmt.Errorf("invalid token"), fiber.StatusUnauthorized)
	}

	if val == "" || val != token {
		return helper.HandleError(c, fmt.Errorf("invalid token"), fiber.StatusUnauthorized)
	}

	c.Locals("email", email)
	c.Locals("user", user)
	c.Locals("token", token)
	return c.Next()
}

func RedirectInterceptor(c *fiber.Ctx) error {
	redirect := c.Query("return_to")
	if redirect != "" {
		c.Locals("redirect", redirect)
	} else {
		return helper.HandleError(c, fmt.Errorf("missing return_to query parameter"), fiber.StatusBadRequest)
	}
	return c.Next()
}

func ClientIdInterceptor(c *fiber.Ctx) error {
	clientId := c.Query("client_id")
	if clientId != "" {
		c.Locals("client_id", clientId)
	} else {
		return helper.HandleError(c, fmt.Errorf("missing client_id query parameter"), fiber.StatusBadRequest)
	}
	return c.Next()
}

func UserIdInterceptor(c *fiber.Ctx, mongoClient *mongo.Client) error {
	userId := c.Query("user_id")
	if userId != "" {
		c.Locals("user_id", userId)
	} else {
		return helper.HandleError(c, fmt.Errorf("missing user_id query parameter"), fiber.StatusBadRequest)
	}

	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return helper.HandleError(c, fmt.Errorf("invalid user_id format"), fiber.StatusBadRequest)
	}

	filter := bson.M{"_id": objectID}
	userCollection := mongoClient.Database("custom-auth").Collection("users")
	var user models.User
	err = userCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return helper.HandleError(c, fmt.Errorf("user not found"), fiber.StatusUnauthorized)
		}
		return helper.HandleError(c, err, fiber.StatusInternalServerError)
	}
	return c.Next()
}
