package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

func init() {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Print("Error loading .env file")
	}

	fmt.Println("Environment variables loaded.")
}
