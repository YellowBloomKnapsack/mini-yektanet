package main

import (
	"log"
	"os"

	"YellowBloomKnapsack/mini-yektanet/eventserver/cache"

	"YellowBloomKnapsack/mini-yektanet/common/grafana"
	"YellowBloomKnapsack/mini-yektanet/common/tokenhandler"
	"YellowBloomKnapsack/mini-yektanet/eventserver/handlers"
	"YellowBloomKnapsack/mini-yektanet/eventserver/producer"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	if err := godotenv.Load(".env", "../common/.env", "../adserver/.env", "../panel/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	tokenHandler := tokenhandler.NewTokenHandlerService()
	redisUrl := os.Getenv("REDIS_HOST")
	redisUrl = redisUrl + ":6379"

	cacheService := cache.NewEventServerCache(redisUrl)
	producerService := producer.NewKafkaProducer()

	handler := handlers.NewEventServerHandler(tokenHandler, cacheService, producerService)

	// Initialize the Gin router
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))

	// Define the routes
	r.Use(grafana.PrometheusMiddleware())

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.POST(os.Getenv("CLICK_REQ_PATH"), handler.PostClick)
	r.POST(os.Getenv("IMPRESSION_REQ_PATH"), handler.PostImpression)

	// Run the server
	port := os.Getenv("EVENT_SERVER_PORT")
	if port == "" {
		port = "8082"
	}
	r.Run(os.Getenv("GIN_HOSTNAME") + ":" + port)
}
