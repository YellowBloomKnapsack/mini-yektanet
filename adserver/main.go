package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"YellowBloomKnapsack/mini-yektanet/adserver/handlers"
	"YellowBloomKnapsack/mini-yektanet/adserver/logic"
)

func main() {
	if err := godotenv.Load(".env", "../common/.env", "../eventserver/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

    logic.Init()
	go logic.StartTicker()

	r := gin.Default()

	r.GET("/:publisherUsername", handlers.GetAd)

	port := os.Getenv("AD_SERVER_PORT")
	if port == "" {
		port = "8081"
	}
	r.Run(":" + port)
}
