package controllers

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/vendz/custom-0auth/models"
	"github.com/vendz/custom-0auth/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (databaseClient Database) SetTokenInRedis(email, rand, token string) error {
	redisKey := "token" + rand + ":" + email
	err := databaseClient.RedisClient.Set(context.Background(), redisKey, token, 100*365*24*time.Hour).Err()
	if err != nil {
		return err
	}
	fmt.Println("Created token in redis for user: ", email)
	return nil
}

func (databaseClient Database) GenerateAndSetToken(email string) (string, error) {
	rand_int := time.Now().Unix()
	rand := strconv.Itoa(int(rand_int))
	token, err := utils.GenerateToken(email, rand, os.Getenv("JWT_SECRET"))
	if err != nil {
		return "", err
	}
	err = databaseClient.SetTokenInRedis(email, rand, token)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (databaseClient Database) GetCurrentUser(email string) (*models.GetUserRequest, error) {
	userCollection := databaseClient.MongoClient.Database("custom-auth").Collection("users")
	filter := bson.M{"email": email}
	user := models.GetUserRequest{}
	err := userCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
