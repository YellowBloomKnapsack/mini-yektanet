package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

    "YellowBloomKnapsack/mini-yektanet/AdServer/handlers"
    "YellowBloomKnapsack/mini-yektanet/AdServer/logic"
)

func main() {
	if err := godotenv.Load(".env", "../common/.env", "../EventServer/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

    go logic.StartTicker()

	r := gin.Default()

	r.GET("/", handlers.GetAd)

	port := os.Getenv("AD_SERVER_PORT")
	if port == "" {
		port = "8081"
	}
	r.Run(":" + port)
}
