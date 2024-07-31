package main

import (
	"html/template"
	"log"
	"os"
	"strconv"

	"YellowBloomKnapsack/mini-yektanet/panel/database"
	panelGrafana "YellowBloomKnapsack/mini-yektanet/panel/grafana"
	"YellowBloomKnapsack/mini-yektanet/panel/handlers"
	"YellowBloomKnapsack/mini-yektanet/panel/reporter"

	"github.com/chenjiandongx/ginprom"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	if err := godotenv.Load(".env", "../common/.env", "../publisherwebsite/.env", "../adserver/.env", "../eventserver/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	database.InitDB()

	clickTopic := os.Getenv("KAFKA_TOPIC_CLICK")
	impressionTopic := os.Getenv("KAFKA_TOPIC_IMPRESSION")
	clickBuffLimit, _ := strconv.Atoi(os.Getenv("KAFKA_CLICK_BUFF_LIMIT"))
	impressionBuffLimit, _ := strconv.Atoi(os.Getenv("KAFKA_IMPRESSION_BUFF_LIMIT"))

	reporterService := reporter.NewReporterService(clickTopic, impressionTopic, clickBuffLimit, impressionBuffLimit)
	reporterService.Start()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))
	r.Use(ginprom.PromMiddleware(nil))
	r.SetFuncMap(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"until": func(count int) []int {
			var i int
			var items []int
			for i = 0; i < count; i++ {
				items = append(items, i)
			}
			return items
		},
	})

	r.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))
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

	port := os.Getenv("PANEL_PORT")
	if port == "" {
		port = "8083"
	}
	panelGrafana.InitializeMetrics()

	r.Run(os.Getenv("GIN_HOSTNAME") + ":" + port)
}
