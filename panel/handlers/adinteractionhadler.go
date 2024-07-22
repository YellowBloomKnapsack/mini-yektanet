package handlers

import (
	"net/http"

	"YellowBloomKnapsack/mini-yektanet/common/dto"
	"YellowBloomKnapsack/mini-yektanet/common/models"
	"YellowBloomKnapsack/mini-yektanet/panel/database"

	"github.com/gin-gonic/gin"
)
func HandleAdInteraction(interactionType models.AdsInteractionType) gin.HandlerFunc {
    return func(c *gin.Context) {
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
            ClickTime:   request.ClickTime,
            AdID:        ad.ID,
            PublisherID: publisher.ID,
        }

        if err := database.DB.Create(&interaction).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create interaction"})
            return
        }

        c.JSON(http.StatusCreated, gin.H{"message": "Interaction recorded successfully"})
    }
}

