package handlers

import (
	"net/http"
	"os"
	"strconv"

	"YellowBloomKnapsack/mini-yektanet/common/dto"
	"YellowBloomKnapsack/mini-yektanet/common/models"
	"YellowBloomKnapsack/mini-yektanet/panel/database"
	"YellowBloomKnapsack/mini-yektanet/panel/logic"

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
	var impressionsCount int64
	for _, ad := range ads {
		_ = database.DB.Model(&models.AdsInteraction{}).Where("ad_id = ? AND type = ?", ad.ID, models.Impression).Count(&impressionsCount)
		adDTO := dto.AdDTO{
			ID:                ad.ID,
			Text:              ad.Text,
			ImagePath:         "http://" + os.Getenv("HOSTNAME") + ":" + os.Getenv("PANEL_PORT") + ad.ImagePath,
			Bid:               ad.Bid,
			Website:           ad.Website,
			TotalCost:         ad.TotalCost,
			BalanceAdvertiser: ad.Advertiser.Balance,
			Impressions:       int(impressionsCount),
			Score:             logic.GetScore(ad, int(impressionsCount)),
		}
		adDTOs = append(adDTOs, adDTO)
	}

	c.JSON(http.StatusOK, adDTOs)

}

func NotifyAdsBrake(adId uint) {
	adIdStr := strconv.FormatUint(uint64(adId), 10)

	notifyApiPath := "http://" + os.Getenv("HOSTNAME") + ":" + os.Getenv("AD_SERVER_PORT") + os.Getenv("NOTIFY_API_PATH") + "/" + adIdStr
	resp, err := http.Post(notifyApiPath, "", nil)
	if err != nil {
		return
	}

	defer resp.Body.Close()
}
