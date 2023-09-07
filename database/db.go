package database

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	controllers "github.com/vendz/custom-0auth/controllers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDatabase() controllers.Database {
	client, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(os.Getenv("MONGO_URI_LOCAL")))
	if err != nil {
		panic(err)
	}

	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("connected to MongoDB...")

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URI_LOCAL"),
		Password: "",
		DB:       0,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}
	fmt.Println("connected to Redis...")

	return controllers.Database{
		MongoClient: client,
		RedisClient: rdb,
	}
}
