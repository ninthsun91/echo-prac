package main

import (
	"log"

	"github.com/joho/godotenv"

	server "myapp/internal"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	server := server.NewServer()
	server.Start(":1323")
}
