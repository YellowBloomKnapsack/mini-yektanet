package main

import (
	"log"
	"os"

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

	port := os.Getenv("PANEL_PORT")
	if port == "" {
		port = "8083"
	}
	r.Run(":" + port)
}
