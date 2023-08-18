package controllers

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vendz/custom-0auth/models"
	"github.com/vendz/custom-0auth/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (databaseClient Database) AuthenticateUser(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "auth": true})
}

func (databaseClient Database) CreateUser(c *fiber.Ctx) error {
	var payload models.CreateUserRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": errors})
	}

	userCollection := databaseClient.MongoClient.Database("custom-auth").Collection("users")

	filter := bson.M{"email": payload.Email}
	count, _ := userCollection.CountDocuments(context.TODO(), filter)
	if count > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "error": "Email already exists"})
	}

	token, err := utils.GenerateToken(*payload.Email, os.Getenv("JWT_SECRET"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	now := time.Now()
	hash, _ := utils.HashPassword(*payload.Password)

	newUser := models.User{
		Id:          primitive.NewObjectID(),
		Name:        payload.Name,
		Email:       payload.Email,
		Token:       &token,
		Password:    &hash,
		PhoneNumber: payload.PhoneNumber,
		Bio:         payload.Bio,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	result, err := userCollection.InsertOne(context.Background(), newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}
	fmt.Println("Created user with id: ", result.InsertedID)

	redisKey := "token:" + *payload.Email
	err = databaseClient.RedisClient.Set(context.Background(), redisKey, token, 100*365*24*time.Hour).Err()
	if err != nil {
		userCollection.DeleteOne(context.Background(), bson.M{"_id": result.InsertedID})
		fmt.Println("Deleting user with id: ", result.InsertedID, " due to redis error: ", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "error": err.Error()})
	} else {
		fmt.Println("Created token in redis for user: ", *payload.Email)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "user": result.InsertedID, "token": token})
}

func (databaseClient Database) LoginUser(c *fiber.Ctx) error {
	var payload models.LoginUserRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	ctx := context.Background()
	userCollection := databaseClient.MongoClient.Database("custom-auth").Collection("users")
	findUser := models.User{}
	filter := bson.M{"email": payload.Email}

	err := userCollection.FindOne(ctx, filter).Decode(&findUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "error": "Invalid credentials"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	if !utils.CheckPasswordHash(*payload.Password, *findUser.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "error": "Invalid credentials"})
	}

	token, err := utils.GenerateToken(*payload.Email, os.Getenv("JWT_SECRET"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	update := bson.M{
		"$set": bson.M{
			"token": token,
		},
	}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	redisKey := "token:" + *payload.Email
	err = databaseClient.RedisClient.Set(context.Background(), redisKey, token, 100*365*24*time.Hour).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "error": err.Error()})
	} else {
		fmt.Println("Created token in redis for user: ", *payload.Email)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "user": findUser.Id, "token": token})
}

func (databaseClient Database) LogoutUser(c *fiber.Ctx) error {
	rdbStr := "token:" + c.Locals("email").(string)
	_, err := databaseClient.RedisClient.Del(context.Background(), rdbStr).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

func (databaseClient Database) GetMe(c *fiber.Ctx) error {
	userCollection := databaseClient.MongoClient.Database("custom-auth").Collection("users")
	findUser := models.User{}
	filter := bson.M{"email": c.Locals("email").(string)}

	err := userCollection.FindOne(context.Background(), filter).Decode(&findUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "error": "Invalid credentials"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "user": findUser})
}
