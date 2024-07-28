package main

import (
	"log"
	"os"

	"YellowBloomKnapsack/mini-yektanet/panel/database"
	"YellowBloomKnapsack/mini-yektanet/panel/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env", "../common/.env", "../publisherwebsite/.env", "../adserver/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	database.InitDB()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	advertiser := r.Group("/advertiser")
	{
		advertiser.GET("/:username/panel", handlers.AdvertiserPanel)
		advertiser.POST("/:username/add-funds", handlers.AddFunds)
		advertiser.POST("/:username/create-ad", handlers.CreateAd)
		advertiser.POST("/:username/toggle-ad", handlers.ToggleAd)
		advertiser.GET("/:username/ad-report/:id", handlers.AdReport)
		advertiser.POST("/:username/edit-ad", handlers.HandleEditAd)

	}

	publisher := r.Group("/publisher")
	{
		publisher.GET("/:username/panel", handlers.PublisherPanel)
		publisher.POST("/:username/withdraw", handlers.WithdrawPublisherBalance)
	}
	r.GET(os.Getenv("GET_ADS_API"), handlers.GetActiveAds)
	r.POST(os.Getenv("INTERACTION_CLICK_API"), handlers.HandleClickAdInteraction)
	r.POST(os.Getenv("INTERACTION_IMPRESSION_API"), handlers.HandleImpressionAdInteraction)

	port := os.Getenv("PANEL_PORT")
	if port == "" {
		port = "8083"
	}
	r.Run(os.Getenv("GIN_HOSTNAME") + ":" + port)
}
