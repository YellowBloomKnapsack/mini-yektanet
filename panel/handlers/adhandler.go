package handlers

import (
	"gorm.io/gorm"
	"net/http"
	"os"
	"strconv"
	"time"

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

func NotifyAdsBrake(adId uint) {
	adIdStr := strconv.FormatUint(uint64(adId), 10)

	notifyApiPath := "http://" + os.Getenv("HOSTNAME") + ":" + os.Getenv("AD_SERVER_PORT") + os.Getenv("NOTIFY_API_PATH") + "/" + adIdStr
	resp, err := http.Post(notifyApiPath, "", nil)
	if err != nil {
		return
	}

	defer resp.Body.Close()
}

func getSumOfBids(db *gorm.DB, adID uint) (int64, error) {
	var sum int64
	eventType := 1
	now := time.Now()
	startTime := now.Add(+(210 * time.Minute))
	endTime := now.Add(-15*time.Second + 210*time.Minute)

	// GORM query to sum bids
	err := db.Model(&models.AdsInteraction{}).
		Where("ad_id = ? AND type = ? AND event_time BETWEEN ? AND ?", adID, eventType, endTime, startTime).
		Select("SUM(bid)").
		Scan(&sum).Error

	if err != nil {
		return 0, err
	}

	return sum, nil
}
