package main

import (
	"log"
	"os"

	// "YellowBloomKnapsack/mini-yektanet/common/database"
	"YellowBloomKnapsack/mini-yektanet/eventserver/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/gin-contrib/cors"
)

func main() {
	if err := godotenv.Load(".env", "../common/.env", "../adserver/.env", "../panel/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()
	
	r.Use(cors.New(cors.Config{
			AllowAllOrigins: true,
		}))

	r.POST(os.Getenv("CLICK_REQ_PATH"), handlers.PostClick)

	r.POST(os.Getenv("IMPRESSION_REQ_PATH"), handlers.PostImpression)

	port := os.Getenv("EVENT_SERVER_PORT")
	if port == "" {
		port = "8082"
	}
	r.Run(":" + port)
}
