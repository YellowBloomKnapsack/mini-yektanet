package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"YellowBloomKnapsack/mini-yektanet/adserver/cache"
	"YellowBloomKnapsack/mini-yektanet/adserver/handlers"
	"YellowBloomKnapsack/mini-yektanet/adserver/logic"
	"YellowBloomKnapsack/mini-yektanet/common/tokenhandler"
)

func main() {
	if err := godotenv.Load(".env", "../common/.env", "../eventserver/.env", "../panel/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	tokenHandler := tokenhandler.NewTokenHandlerService()
	redisUrl := os.Getenv("REDIS_URL")
	cacheService := cache.NewAdServerCache(redisUrl)
	logicService := logic.NewLogicService(cacheService)

	handler := handlers.NewAdServerHandler(logicService, tokenHandler)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))

	r.GET("/:publisherId", handler.GetAd)
	r.POST(os.Getenv("NOTIFY_API_PATH")+"/:adId", handler.BrakeAd)

	port := os.Getenv("AD_SERVER_PORT")
	if port == "" {
		port = "8081"
	}
	r.Run(os.Getenv("GIN_HOSTNAME") + ":" + port)
}
