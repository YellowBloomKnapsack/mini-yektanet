package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"YellowBloomKnapsack/mini-yektanet/common/tokenhandler"
	"YellowBloomKnapsack/mini-yektanet/eventserver/cache"
	"YellowBloomKnapsack/mini-yektanet/eventserver/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env", "../common/.env", "../adserver/.env", "../panel/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	tokenHandler := tokenhandler.NewTokenHandlerService()

	// Initialize the cache service
	redisUrl := os.Getenv("REDIS_URL")
	redisExpireHourStr := os.Getenv("REDIS_EXPIRE_HOUR")
	redisExpireHour, err := strconv.Atoi(redisExpireHourStr)
	if err != nil {
		log.Fatalf("Invalid REDIS_EXPIRE_HOUR value: %v", err)
	}
	cacheService := cache.NewCacheService(time.Duration(redisExpireHour)*time.Hour, redisUrl)

	// Initialize the event server handler
	handler := handlers.NewEventServerHandler(tokenHandler, cacheService)

	// Initialize the Gin router
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))

	// Define the routes
	r.POST(os.Getenv("CLICK_REQ_PATH"), handler.PostClick)
	r.POST(os.Getenv("IMPRESSION_REQ_PATH"), handler.PostImpression)

	// Run the server
	port := os.Getenv("EVENT_SERVER_PORT")
	if port == "" {
		port = "8082"
	}
	r.Run(os.Getenv("EVENT_SERVER_HOSTNAME") + ":" + port)
}
