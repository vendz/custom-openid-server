package controllers

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/vendz/custom-0auth/utils"
)

func (databaseClient Database) setTokenInRedis(email, token string) error {
	redisKey := "token:" + email
	err := databaseClient.RedisClient.Set(context.Background(), redisKey, token, 100*365*24*time.Hour).Err()
	if err != nil {
		return err
	}
	fmt.Println("Created token in redis for user: ", email)
	return nil
}

func (databaseClient Database) generateAndSetToken(email string) (string, error) {
	token, err := utils.GenerateToken(email, os.Getenv("JWT_SECRET"))
	if err != nil {
		return "", err
	}
	err = databaseClient.setTokenInRedis(email, token)
	if err != nil {
		return "", err
	}
	return token, nil
}
