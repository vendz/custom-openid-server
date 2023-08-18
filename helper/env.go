package helper

import (
	"fmt"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Enviroment variables loaded...")
}
