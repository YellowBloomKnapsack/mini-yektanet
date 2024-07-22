package main

import (
	"log"
	"os"

	// "YellowBloomKnapsack/mini-yektanet/common/database"
	"YellowBloomKnapsack/mini-yektanet/eventserver/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env", "../common/.env", "../adserver/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	r.POST(os.Getenv("CLICK_REQ_PATH"), handlers.PostClick)

	port := os.Getenv("EVENT_SERVER_PORT")
	if port == "" {
		port = "8082"
	}
	r.Run(":" + port)
}
