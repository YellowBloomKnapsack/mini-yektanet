package handlers

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"YellowBloomKnapsack/mini-yektanet/common/dto"
	"YellowBloomKnapsack/mini-yektanet/common/models"
	"YellowBloomKnapsack/mini-yektanet/panel/database"
	"YellowBloomKnapsack/mini-yektanet/panel/logic"

	"github.com/gin-gonic/gin"
)

func GetActiveAds(c *gin.Context) {
	var ads []models.Ad
	result := database.DB.Preload("Advertiser").
		Preload("AdsInteractions").
		Preload("Keywords").
		Joins("JOIN advertisers ON advertisers.id = ads.advertiser_id").
		Where("ads.active = ? AND advertisers.balance > ?", true, 0).
		Find(&ads)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ads"})
		return
	}

	var adDTOs []dto.AdDTO
	for _, ad := range ads {
		impressionsCount := getImpressionCounts(ad.AdsInteractions)

		var keywordStrings []string
		for _, keyword := range ad.Keywords {
			cleanKeyword := strings.TrimSpace(keyword.Keywords)
			if cleanKeyword != "" {
				keywordStrings = append(keywordStrings, cleanKeyword)
			}
		}

		// Join keywords into a comma-separated string
		keywordString := strings.Join(keywordStrings, ", ")

		adDTO := dto.AdDTO{
			ID:                ad.ID,
			Text:              ad.Text,
			ImagePath:         "http://" + os.Getenv("PANEL_PUBLIC_HOSTNAME") + ":" + os.Getenv("PANEL_PORT") + ad.ImagePath,
			Bid:               ad.Bid,
			Website:           ad.Website,
			TotalCost:         ad.TotalCost,
			BalanceAdvertiser: ad.Advertiser.Balance,
			Impressions:       impressionsCount,
			Score:             logic.GetScore(ad, impressionsCount),
			Keywords:          keywordString,
		}
		adDTOs = append(adDTOs, adDTO)
	}

	c.JSON(http.StatusOK, adDTOs)

}

func NotifyAdsBrake(adId uint) {
	adIdStr := strconv.FormatUint(uint64(adId), 10)

	notifyApiPath := "http://" + os.Getenv("AD_SERVER_HOSTNAME") + ":" + os.Getenv("AD_SERVER_PORT") + os.Getenv("NOTIFY_API_PATH") + "/" + adIdStr
	resp, err := http.Post(notifyApiPath, "", nil)
	if err != nil {
		return
	}

	defer resp.Body.Close()
}

func getImpressionCounts(adsInteractions []models.AdsInteraction) int {
	count := 0
	for _, interaction := range adsInteractions {
		if interaction.Type == int(models.Impression) {
			count++
		}
	}
	return count
}
