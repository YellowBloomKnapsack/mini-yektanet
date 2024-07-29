package logic

import (
	"database/sql"
	"os"
	"strconv"
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

	if sum.Valid {
		return int64(sum.Int64), nil
	} else {
		return 0, err
	}
}
