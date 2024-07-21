package main

import (
	"log"
	"os"

	// "YellowBloomKnapsack/mini-yektanet/common/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env", "../common/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	// TODO: convert to api call
	// database.InitDB()

	r := gin.Default()

	port := os.Getenv("EVENT_SERVER_PORT")
	if port == "" {
		port = "8082"
	}
	r.Run(":" + port)
}
