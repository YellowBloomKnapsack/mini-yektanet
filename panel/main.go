package panel

import (
	"log"
	"os"

	"YellowBloomKnapsack/mini-yektanet/panel/handlers"
	"YellowBloomKnapsack/mini-yektanet/common/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env", "../common/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	database.InitDB()

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/advertiser/:username/panel", handlers.AdvertiserPanel)
	r.POST("/advertiser/:username/add-funds", handlers.AddFunds)
	r.POST("/advertiser/:username/create-ad", handlers.CreateAd)
	r.POST("/advertiser/:username/toggle-ad", handlers.ToggleAd)
	r.GET("/advertiser/:username/ad-report/:id", handlers.AdReport)

	port := os.Getenv("PANEL_PORT")
	if port == "" {
		port = "8083"
	}
	r.Run(":" + port)
}
