package handlers

import (
	"net/http"

	"YellowBloomKnapsack/mini-yektanet/common/dto"
	"YellowBloomKnapsack/mini-yektanet/common/models"
	"YellowBloomKnapsack/mini-yektanet/panel/database"

	"github.com/gin-gonic/gin"
)

func GetActiveAds(c *gin.Context) {

	var ads []models.Ad
	result := database.DB.Where("active = ?", true).Find(&ads)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ads"})
		return
	}

	var adDTOs []dto.AdDTO
	for _, ad := range ads {
		adDTO := dto.AdDTO{
			ID:        ad.ID,
			Text:      ad.Text,
			ImagePath: ad.ImagePath,
			Bid:       ad.Bid,
			Website:   "", // You might want to add a Website field to your Ad model if needed
		}
		adDTOs = append(adDTOs, adDTO)
	}

	c.JSON(http.StatusOK, adDTOs)

}
