package main

import (
	"YellowBloomKnapsack/mini-yektanet/common/grafana"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var currentSiteNames []string

func main() {
	if err := godotenv.Load(".env", "../common/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	currentSiteNames = []string{
		"varzesh3",
		"zoomit",
		"digikala",
		"sheypoor",
		"filimo",
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))
	r.Use(grafana.PrometheusMiddleware())

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/:siteName", getSite)
	r.LoadHTMLGlob("html/*")
	r.Static("/static", "./static")

	port := os.Getenv("PUBLISHER_WEBSITE_PORT")
	if port == "" {
		port = "8084"
	}
	r.Run(os.Getenv("GIN_HOSTNAME") + ":" + port)
}

func getSite(c *gin.Context) {
	siteName := c.Param("siteName")
	for _, name := range currentSiteNames {
		if siteName == name {
			c.HTML(http.StatusOK, siteName+".html", gin.H{
				"title": siteName,
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "not here"})
}
