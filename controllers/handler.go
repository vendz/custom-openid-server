package controllers

import (
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Database struct {
	MongoClient *mongo.Client
	RedisClient *redis.Client
}
