package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"YellowBloomKnapsack/mini-yektanet/common/tokenhandler"
	"YellowBloomKnapsack/mini-yektanet/eventserver/cache"
	"YellowBloomKnapsack/mini-yektanet/eventserver/handlers"
	"YellowBloomKnapsack/mini-yektanet/eventserver/worker"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env", "../common/.env", "../adserver/.env", "../panel/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	tokenHandler := tokenhandler.NewTokenHandlerService()
	worker := worker.NewWorkerService()

	redisUrl := os.Getenv("REDIS_URL")
	redisExpireHour, _ := strconv.Atoi("REDIS_EXPIRE_HOUR")

	cache := cache.NewCacheService(time.Duration(redisExpireHour)*time.Hour, redisUrl)

	handler := handlers.NewEventServerHandler(tokenHandler, worker, cache)
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))

	r.POST(os.Getenv("CLICK_REQ_PATH"), handler.PostClick)

	r.POST(os.Getenv("IMPRESSION_REQ_PATH"), handler.PostImpression)

	port := os.Getenv("EVENT_SERVER_PORT")
	if port == "" {
		port = "8082"
	}
	r.Run(os.Getenv("EVENT_SERVER_HOSTNAME") + ":" + port)
}
