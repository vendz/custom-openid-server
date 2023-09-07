package controllers

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vendz/custom-0auth/helper"
	"github.com/vendz/custom-0auth/models"
	"github.com/vendz/custom-0auth/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var backBroundCtx = context.Background()

func (databaseClient Database) SingleSignon(c *fiber.Ctx) error {
	userCollection := databaseClient.MongoClient.Database("custom-auth").Collection("users")
	findUser := models.User{}
	filter := bson.M{"email": c.Locals("email").(string)}

	err := userCollection.FindOne(backBroundCtx, filter).Decode(&findUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return helper.HandleError(c, err, fiber.StatusBadRequest)
		}
		return helper.HandleError(c, err, fiber.StatusInternalServerError)
	}

	clientAccessToken, err := utils.GenerateToken(c.Locals("email").(string), c.Locals("client_id").(string))
	if err != nil {
		return helper.HandleError(c, err, fiber.StatusInternalServerError)
	}

	tempMap := findUser.AccessTokens
	tempMap[c.Locals("client_id").(string)] = clientAccessToken
	update := bson.M{
		"$set": bson.M{
			"accesstokens": tempMap,
		},
	}

	_, err = userCollection.UpdateOne(backBroundCtx, filter, update)
	if err != nil {
		return helper.HandleError(c, err, fiber.StatusInternalServerError)
	}

	redisKey := c.Locals("client_id").(string) + ":" + c.Locals("email").(string)
	err = databaseClient.RedisClient.Set(backBroundCtx, redisKey, clientAccessToken, time.Hour).Err()
	if err != nil {
		return helper.HandleError(c, err, fiber.StatusInternalServerError)
	}

	userData := map[string]string{
		"name":        *findUser.Name,
		"email":       *findUser.Email,
		"phoneNumber": *findUser.PhoneNumber,
		"bio":         *findUser.Bio,
	}

	return helper.FormatResponse(c, fiber.StatusOK, "success", fiber.Map{"user": userData, "token": findUser.Token, "client_token": clientAccessToken, "return_to": c.Locals("redirect").(string)})
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

	fmt.Println(payload.Name, payload.Email, payload.Password, payload.PhoneNumber, payload.Bio)

	userCollection := databaseClient.MongoClient.Database("custom-auth").Collection("users")

	filter := bson.M{"email": payload.Email}
	count, _ := userCollection.CountDocuments(backBroundCtx, filter)
	if count > 0 {
		return helper.HandleError(c, fmt.Errorf("email already exists"), fiber.StatusBadRequest)
	}

	token, err := utils.GenerateToken(*payload.Email, os.Getenv("JWT_SECRET"))
	if err != nil {
		return helper.HandleError(c, err, fiber.StatusInternalServerError)
	}

	clientAccessToken, err := utils.GenerateToken(*payload.Email, c.Locals("client_id").(string))
	if err != nil {
		return helper.HandleError(c, err, fiber.StatusInternalServerError)
	}
	tokenMap := map[string]string{
		c.Locals("client_id").(string): clientAccessToken,
	}

	now := time.Now()
	hash, _ := utils.HashPassword(*payload.Password)

	newUser := models.User{
		Id:           primitive.NewObjectID(),
		Name:         payload.Name,
		Email:        payload.Email,
		Token:        &token,
		AccessTokens: tokenMap,
		Password:     &hash,
		PhoneNumber:  payload.PhoneNumber,
		Bio:          payload.Bio,
		IsClient:     false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	result, err := userCollection.InsertOne(backBroundCtx, newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}
	fmt.Println("Created user with id: ", result.InsertedID)

	redisKey := "token:" + *payload.Email
	redisKey2 := c.Locals("client_id").(string) + ":" + *payload.Email

	err = databaseClient.RedisClient.Set(backBroundCtx, redisKey, token, 100*365*24*time.Hour).Err()
	err2 := databaseClient.RedisClient.Set(backBroundCtx, redisKey2, clientAccessToken, time.Hour).Err()

	if err != nil && err2 != nil {
		userCollection.DeleteOne(backBroundCtx, bson.M{"_id": result.InsertedID})
		fmt.Println("Deleting user with id: ", result.InsertedID, " due to redis error: ", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "error": err.Error()})
	} else {
		fmt.Println("Created token in redis for user: ", *payload.Email)
	}

	userData := map[string]string{
		"id":          result.InsertedID.(primitive.ObjectID).Hex(),
		"name":        *payload.Name,
		"email":       *payload.Email,
		"phoneNumber": *payload.PhoneNumber,
		"bio":         *payload.Bio,
	}

	return helper.FormatResponse(c, fiber.StatusOK, "success", fiber.Map{"user": userData, "token": token, "access_token": clientAccessToken, "return_to": c.Locals("redirect").(string)})
}

func (databaseClient Database) LoginUser(c *fiber.Ctx) error {
	var payload models.LoginUserRequest
	if err := c.BodyParser(&payload); err != nil {
		return helper.HandleError(c, err, fiber.StatusBadRequest)
	}

	ctx := backBroundCtx
	userCollection := databaseClient.MongoClient.Database("custom-auth").Collection("users")
	findUser := models.User{}
	filter := bson.M{"email": payload.Email}

	err := userCollection.FindOne(ctx, filter).Decode(&findUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return helper.HandleError(c, fmt.Errorf("invalid credentials"), fiber.StatusBadRequest)
		}
		return helper.HandleError(c, err, fiber.StatusInternalServerError)
	}

	if !utils.CheckPasswordHash(*payload.Password, *findUser.Password) {
		return helper.HandleError(c, fmt.Errorf("invalid credentials"), fiber.StatusBadRequest)
	}

	token, err := utils.GenerateToken(*payload.Email, os.Getenv("JWT_SECRET"))
	if err != nil {
		return helper.HandleError(c, err, fiber.StatusInternalServerError)
	}

	clientAccessToken, err := utils.GenerateToken(*payload.Email, c.Locals("client_id").(string))
	if err != nil {
		return helper.HandleError(c, err, fiber.StatusInternalServerError)
	}

	tempMap := findUser.AccessTokens
	tempMap[c.Locals("client_id").(string)] = clientAccessToken
	update := bson.M{
		"$set": bson.M{
			"token":        token,
			"accesstokens": tempMap,
		},
	}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return helper.HandleError(c, err, fiber.StatusInternalServerError)
	}

	redisKey := "token:" + *payload.Email
	redisKey2 := c.Locals("client_id").(string) + ":" + *payload.Email

	err = databaseClient.RedisClient.Set(backBroundCtx, redisKey, token, 100*365*24*time.Hour).Err()
	err2 := databaseClient.RedisClient.Set(backBroundCtx, redisKey2, clientAccessToken, time.Hour).Err()

	if err != nil && err2 != nil {
		return helper.HandleError(c, err, fiber.StatusInternalServerError)
	} else {
		fmt.Println("Created token in redis for user: ", *payload.Email)
	}

	userData := map[string]string{
		"id":          findUser.Id.Hex(),
		"name":        *findUser.Name,
		"email":       *findUser.Email,
		"phoneNumber": *findUser.PhoneNumber,
		"bio":         *findUser.Bio,
	}

	return helper.FormatResponse(c, fiber.StatusOK, "success", fiber.Map{"user": userData, "token": token, "client_token": clientAccessToken, "return_to": c.Locals("redirect").(string)})
}

func (databaseClient Database) LogoutUser(c *fiber.Ctx) error {
	rdbStr := "token:" + c.Locals("email").(string)
	_, err := databaseClient.RedisClient.Del(backBroundCtx, rdbStr).Result()
	if err != nil {
		return helper.HandleError(c, err, fiber.StatusInternalServerError)
	}

	return helper.FormatResponse(c, fiber.StatusOK, "success", fiber.Map{"message": "logged out"})
}

func (databaseClient Database) GetMe(c *fiber.Ctx) error {
	userCollection := databaseClient.MongoClient.Database("custom-auth").Collection("users")
	findUser := models.User{}
	filter := bson.M{"email": c.Locals("email").(string)}

	err := userCollection.FindOne(backBroundCtx, filter).Decode(&findUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return helper.HandleError(c, fmt.Errorf("invalid credentials"), fiber.StatusBadRequest)
		}
		return helper.HandleError(c, err, fiber.StatusInternalServerError)
	}

	userData := map[string]string{
		"name":        *findUser.Name,
		"email":       *findUser.Email,
		"phoneNumber": *findUser.PhoneNumber,
		"bio":         *findUser.Bio,
	}

	return helper.FormatResponse(c, fiber.StatusOK, "success", fiber.Map{"user": userData})
}

func (databaseClient Database) AuthenticateUser(c *fiber.Ctx) error {
	return helper.FormatResponse(c, fiber.StatusOK, "success", fiber.Map{"message": "authenticated"})
}
