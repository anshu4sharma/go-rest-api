package main

import (
	"log"
	"os"

	"github.com/anshu4sharma/go-rest-api/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize and start the application
	application := app.NewApp()
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	if err := application.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}
