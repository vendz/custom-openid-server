package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vendz/custom-0auth/helper"
	"github.com/vendz/custom-0auth/models"
	"github.com/vendz/custom-0auth/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (databaseClient Database) SingleSignon(c *fiber.Ctx) error {
	return helper.FormatResponse(c, fiber.StatusOK, "success", fiber.Map{"user": c.Locals("user"), "token": c.Locals("token"), "return_to": c.Locals("redirect").(string)})
}

func (databaseClient Database) CreateUser(c *fiber.Ctx) error {
	var backBroundCtx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var payload models.CreateUserRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": errors})
	}

	fmt.Println(payload.Name, payload.Email, payload.Password, payload.PhoneNumber, payload.Bio)

	userCollection := databaseClient.MongoClient.Database("custom-auth").Collection("users")

	filter := bson.M{"email": payload.Email}
	count, _ := userCollection.CountDocuments(backBroundCtx, filter)

	if count > 0 {
		return helper.HandleError(c, fmt.Errorf("email already exists"), fiber.StatusBadRequest)
	}

	token, err := databaseClient.GenerateAndSetToken(*payload.Email)
	if err != nil {
		return helper.HandleError(c, err, fiber.StatusInternalServerError)
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
		IsClient:    false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	result, err := userCollection.InsertOne(backBroundCtx, newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}
	fmt.Println("Created user with id: ", result.InsertedID)

	userData := map[string]string{
		"id":          result.InsertedID.(primitive.ObjectID).Hex(),
		"name":        *payload.Name,
		"email":       *payload.Email,
		"phoneNumber": *payload.PhoneNumber,
		"bio":         *payload.Bio,
	}

	return helper.FormatResponse(c, fiber.StatusOK, "success", fiber.Map{"user": userData, "access_token": token, "return_to": c.Locals("redirect").(string)})
}

func (databaseClient Database) LoginUser(c *fiber.Ctx) error {
	var backBroundCtx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var payload models.LoginUserRequest
	if err := c.BodyParser(&payload); err != nil {
		return helper.HandleError(c, err, fiber.StatusBadRequest)
	}

	userCollection := databaseClient.MongoClient.Database("custom-auth").Collection("users")
	findUser := models.User{}
	filter := bson.M{"email": payload.Email}

	err := userCollection.FindOne(backBroundCtx, filter).Decode(&findUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return helper.HandleError(c, fmt.Errorf("invalid credentials"), fiber.StatusBadRequest)
		}
		return helper.HandleError(c, err, fiber.StatusInternalServerError)
	}

	if !utils.CheckPasswordHash(*payload.Password, *findUser.Password) {
		return helper.HandleError(c, fmt.Errorf("invalid credentials"), fiber.StatusBadRequest)
	}

	token, err := databaseClient.GenerateAndSetToken(*payload.Email)
	if err != nil {
		return helper.HandleError(c, err, fiber.StatusInternalServerError)
	}

	userData := map[string]string{
		"id":          findUser.Id.Hex(),
		"name":        *findUser.Name,
		"email":       *findUser.Email,
		"phoneNumber": *findUser.PhoneNumber,
		"bio":         *findUser.Bio,
	}

	return helper.FormatResponse(c, fiber.StatusOK, "success", fiber.Map{"user": userData, "access_token": token, "return_to": c.Locals("redirect").(string)})
}

func (databaseClient Database) LogoutUser(c *fiber.Ctx) error {
	var backBroundCtx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	rdbStr := "token" + c.Locals("rand").(string) + ":" + c.Locals("email").(string)
	_, err := databaseClient.RedisClient.Del(backBroundCtx, rdbStr).Result()
	defer cancel()
	if err != nil {
		return helper.HandleError(c, err, fiber.StatusInternalServerError)
	}

	return helper.FormatResponse(c, fiber.StatusOK, "success", fiber.Map{"message": "logged out"})
}

func (databaseClient Database) GetMe(c *fiber.Ctx) error {
	return helper.FormatResponse(c, fiber.StatusOK, "success", fiber.Map{"user": c.Locals("user")})
}
