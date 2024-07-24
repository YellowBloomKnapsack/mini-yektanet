package handlers

import (
	"net/http"
	"os"

	"YellowBloomKnapsack/mini-yektanet/common/dto"
	"YellowBloomKnapsack/mini-yektanet/common/models"
	"YellowBloomKnapsack/mini-yektanet/panel/database"

	"github.com/gin-gonic/gin"
)

func GetActiveAds(c *gin.Context) {

	var ads []models.Ad
	result := database.DB.Preload("Advertiser").Preload("AdsInteractions").Where("active = ?", true).Find(&ads)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ads"})
		return
	}

	var adDTOs []dto.AdDTO
	for _, ad := range ads {
		adDTO := dto.AdDTO{
			ID:                ad.ID,
			Text:              ad.Text,
			ImagePath:         "http://" + os.Getenv("HOSTNAME") + ":" + os.Getenv("PANEL_PORT") + ad.ImagePath,
			Bid:               ad.Bid,
			Website:           ad.Website,
			TotalCost:         ad.TotalCost,
			BalanceAdvertiser: ad.Advertiser.Balance,
			Impressions: func() int {
				count := 0
				for _, interaction := range ad.AdsInteractions {
					if interaction.Type == int(models.Impression) {
						count++
					}
				}
				return count
			}(),
		}
		adDTOs = append(adDTOs, adDTO)
	}

	c.JSON(http.StatusOK, adDTOs)

}
func NotifyAdsUpdate() {
	notifyApiPath := "http://" + os.Getenv("HOSTNAME") + ":" + os.Getenv("AD_SERVER_PORT") + os.Getenv("NOTIFY_API_PATH")
	resp, err := http.Post(notifyApiPath, "", nil)
	if err != nil {
		return
	}

	defer resp.Body.Close()
}
