package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"YellowBloomKnapsack/mini-yektanet/adserver/handlers"
	"YellowBloomKnapsack/mini-yektanet/adserver/logic"
	"github.com/gin-contrib/cors"
)

func main() {
	if err := godotenv.Load(".env", "../common/.env", "../eventserver/.env", "../panel/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	logic.Init()
	go logic.StartTicker()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))
	r.GET("/:publisherUsername", handlers.GetAd)
	// r.Static("/static", "../publisherwebsite/static")

	port := os.Getenv("AD_SERVER_PORT")
	if port == "" {
		port = "8081"
	}
	r.Run(":" + port)
}
