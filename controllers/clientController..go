package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vendz/custom-0auth/helper"
)

func (databaseClient Database) CreateClient(c *fiber.Ctx) error {
	return helper.FormatResponse(c, fiber.StatusOK, "success", fiber.Map{"message": "client created!"})
}

// func (databaseClient Database) CreateClient(c *fiber.Ctx) error {
// 	userCollection := databaseClient.MongoClient.Database("custom-auth").Collection("users")
// 	clientCollection := databaseClient.MongoClient.Database("custom-auth").Collection("clients")
// 	findClient := models.Client{}
// 	findUser := models.User{}

// 	var payload models.CreateUserRequest
// 	if err := c.BodyParser(&payload); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "error": err.Error()})
// 	}

// 	errors := utils.ValidateStruct(payload)
// 	if errors != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": errors})
// 	}

// 	fmt.Println(payload.Name, payload.Email, payload.Password, payload.PhoneNumber, payload.Bio)

// 	userCollection := databaseClient.MongoClient.Database("custom-auth").Collection("users")

// 	filter := bson.M{"email": payload.Email}
// 	count, _ := userCollection.CountDocuments(backBroundCtx, filter)
// 	if count > 0 {
// 		return handleError(c, fmt.Errorf("email already exists"), fiber.StatusBadRequest)
// 	}

// 	token, err := utils.GenerateToken(*payload.Email, os.Getenv("JWT_SECRET"))
// 	if err != nil {
// 		return handleError(c, err, fiber.StatusInternalServerError)
// 	}

// 	clientAccessToken, err := utils.GenerateToken(*payload.Email, c.Locals("client_id").(string))
// 	if err != nil {
// 		return handleError(c, err, fiber.StatusInternalServerError)
// 	}
// 	tokenMap := map[string]string{
// 		c.Locals("client_id").(string): clientAccessToken,
// 	}

// 	now := time.Now()
// 	hash, _ := utils.HashPassword(*payload.Password)

// 	newUser := models.User{
// 		Id:           primitive.NewObjectID(),
// 		Name:         payload.Name,
// 		Email:        payload.Email,
// 		Token:        &token,
// 		AccessTokens: tokenMap,
// 		Password:     &hash,
// 		PhoneNumber:  payload.PhoneNumber,
// 		Bio:          payload.Bio,
// 		CreatedAt:    now,
// 		UpdatedAt:    now,
// 	}

// 	result, err := userCollection.InsertOne(backBroundCtx, newUser)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "error": err.Error()})
// 	}
// 	fmt.Println("Created user with id: ", result.InsertedID)

// 	redisKey := "token:" + *payload.Email
// 	redisKey2 := c.Locals("client_id").(string) + ":" + *payload.Email

// 	err = databaseClient.RedisClient.Set(backBroundCtx, redisKey, token, 100*365*24*time.Hour).Err()
// 	err2 := databaseClient.RedisClient.Set(backBroundCtx, redisKey2, clientAccessToken, time.Hour).Err()

// 	if err != nil && err2 != nil {
// 		userCollection.DeleteOne(backBroundCtx, bson.M{"_id": result.InsertedID})
// 		fmt.Println("Deleting user with id: ", result.InsertedID, " due to redis error: ", err.Error())
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "error": err.Error()})
// 	} else {
// 		fmt.Println("Created token in redis for user: ", *payload.Email)
// 	}

// 	userData := map[string]string{
// 		"name":        *payload.Name,
// 		"email":       *payload.Email,
// 		"phoneNumber": *payload.PhoneNumber,
// 		"bio":         *payload.Bio,
// 	}

// 	return formatResponse(c, fiber.StatusOK, "success", fiber.Map{"user": userData, "token": token, "access_token": clientAccessToken, "return_to": c.Locals("redirect").(string)})
// }
