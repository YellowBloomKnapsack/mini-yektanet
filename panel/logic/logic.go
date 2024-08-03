package logic

import (
	"database/sql"
	"os"
	"strconv"
	"strings"
	"time"

	"YellowBloomKnapsack/mini-yektanet/common/models"

	"gorm.io/gorm"
)

func GetScore(ad models.Ad, impressionsCount int) float64 {
	return float64(ad.TotalCost) / float64(impressionsCount+1)
}

// TODO: check time
func GetSumOfBids(db *gorm.DB, adID uint) (int64, error) {
	var sum sql.NullInt64
	eventType := 1
	now := time.Now()
	startTime := now
	offsetSecs, _ := strconv.Atoi(os.Getenv("AD_COST_CHECK_DURATION_SECS"))
	endTime := now.Add(-time.Duration(offsetSecs) * time.Second)

	// GORM query to sum bids
	err := db.Model(&models.AdsInteraction{}).
		Where("ad_id = ? AND type = ? AND event_time BETWEEN ? AND ?", adID, eventType, endTime, startTime).
		Select("SUM(bid)").
		Debug().
		Scan(&sum).Error

	// return 0 if no row was found
	if sum.Valid {
		return int64(sum.Int64), nil
	} else {
		return 0, err
	}
}

func SplitAndClean(input string) []string {
	// Split the input string by commas
	parts := strings.Split(input, ",")

	// Create a slice to hold the cleaned parts
	var result []string

	// Loop through the parts
	for _, part := range parts {
		// Trim any leading or trailing whitespace
		trimmed := strings.TrimSpace(part)

		// Add the non-empty trimmed part to the result
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
