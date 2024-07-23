package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"YellowBloomKnapsack/mini-yektanet/common/dto"
	"YellowBloomKnapsack/mini-yektanet/common/models"
	"YellowBloomKnapsack/mini-yektanet/panel/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleClickAdInteraction(c *gin.Context) {
	interactionType := models.Click
	var request dto.InteractionDto
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start a transaction
	tx := database.DB.Begin()

	// Find the publisher
	var publisher models.Publisher
	if err := tx.Where("username = ?", request.PublisherUsername).First(&publisher).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Publisher not found"})
		return
	}

	// Find the ad and its associated advertiser
	var ad models.Ad
	if err := tx.Preload("Advertiser").First(&ad, request.AdID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Ad not found"})
		return
	}

	// Create the interaction
	interaction := models.AdsInteraction{
		Type:        int(interactionType),
		EventTime:   request.EventTime,
		AdID:        ad.ID,
		PublisherID: publisher.ID,
	}

	if err := tx.Create(&interaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create interaction"})
		return
	}

	// Update ad's total cost
	if err := tx.Model(&ad).Update("total_cost", gorm.Expr("total_cost + ?", ad.Bid)).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ad's total cost"})
		return
	}

	yektanetPortionString := os.Getenv("YEKTANET_PORTION")

	// Convert the value to an integer
	yektanetPortion, err := strconv.Atoi(yektanetPortionString)
	if err != nil || yektanetPortion < 0 || yektanetPortion > 100 {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error parsing YEKTANET_PORTION environment variable: %v\n", err)})
		return
	}

	// Increase publisher's balance
	publisherPortion := ad.Bid * int64(100-yektanetPortion) / 100
	if err := tx.Model(&publisher).Update("balance", gorm.Expr("balance + ?", publisherPortion)).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update publisher's balance"})
		return
	}

	// Create a new transaction record for advertiser
	transaction_publisher := models.Transaction{
		CustomerID:   publisher.ID,
		CustomerType: models.Customer_Publisher,
		Amount:       publisherPortion,
		Income:       true,
		Successful:   true,
		Time:         time.Now(),
		Description:  "click on ad",
	}
	if err := tx.Create(&transaction_publisher).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	// Decrease advertiser's balance
	if err := tx.Model(&ad.Advertiser).Update("balance", gorm.Expr("balance - ?", ad.Bid)).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update advertiser's balance"})
		return
	}
	// Create a new transaction record for advertiser
	transaction_advertiser := models.Transaction{
		CustomerID:   ad.ID,
		CustomerType: models.Customer_Advertiser,
		Amount:       ad.Bid,
		Income:       false,
		Successful:   true,
		Time:         time.Now(),
		Description:  "click on ad",
	}
	if err := tx.Create(&transaction_advertiser).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Interaction recorded successfully"})

}
func HandleImpressionAdInteraction(c *gin.Context) {
	interactionType := models.Impression
	var request dto.InteractionDto
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the publisher
	var publisher models.Publisher
	if err := database.DB.Where("username = ?", request.PublisherUsername).First(&publisher).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Publisher not found"})
		return
	}

	// Find the ad
	var ad models.Ad
	if err := database.DB.First(&ad, request.AdID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ad not found"})
		return
	}

	// Create the interaction
	interaction := models.AdsInteraction{
		Type:        int(interactionType),
		EventTime:   request.EventTime,
		AdID:        ad.ID,
		PublisherID: publisher.ID,
	}

	if err := database.DB.Create(&interaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create interaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Interaction recorded successfully"})
}
