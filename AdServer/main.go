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

	port := os.Getenv("AD_SERVER_PORT")
	if port == "" {
		port = "8081"
	}
	r.Run(":" + port)
}
