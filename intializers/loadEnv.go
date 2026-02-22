package intializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		// .env is optional—on Railway etc. env vars are set directly
		log.Println("No .env file found, using environment variables")
	}
}
