package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env", "../common/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	port := os.Getenv("PUBLISHER_WEBSITE_PORT")
	if port == "" {
		port = "8084"
	}
	r.Run(":" + port)
}
