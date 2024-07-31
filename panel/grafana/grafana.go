package grafana

import (
        "log"

        "YellowBloomKnapsack/mini-yektanet/common/models"
        "YellowBloomKnapsack/mini-yektanet/panel/database"
        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ActiveAdsCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "active_ads_count",
		Help: "The total number of active ads",
	})

        AdvertisersCount = promauto.NewCounter(prometheus.CounterOpts{
                Name: "advertisers_count",
                Help: "The total number of advertisers",
        })

        PublishersCount = promauto.NewCounter(prometheus.CounterOpts{
                Name: "publishers_count",
                Help: "The total number of publishers",
        })

        ImpressionCount = promauto.NewCounter(prometheus.CounterOpts{
                Name: "impression_count",
                Help: "The total number of impressions",
        })

        ClickCount = promauto.NewCounter(prometheus.CounterOpts{
                Name: "click_count",
                Help: "The total number of clicks",
        })

        TotalRevenue = promauto.NewCounter(prometheus.CounterOpts{
                Name: "total_revenue",
                Help: "The total revenue generated",
        })

        NumberOfBids = promauto.NewCounter(prometheus.CounterOpts{
                Name: "number_bids",
                Help: "total number of bids",
        })

	AverageBid = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "average_bid",
		Help: "The average bid amount",
	})

	TransactionCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "transaction_count",
		Help: "The number of transactions",
	}, []string{"status"})

	TotalAdvertiserBalance = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "total_advertiser_balance",
		Help: "The total balance of all advertisers",
	})

	TotalPublisherBalance = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "total_publisher_balance",
		Help: "The total balance of all publishers",
	})
)

func InitializeMetrics() {
// Initialize ActiveAdsCount
        var activeAdsCount int64
        if err := database.DB.Model(&models.Ad{}).Where("active = ?", true).Count(&activeAdsCount).Error; err != nil {
                log.Printf("Error counting active ads: %v", err)
        }
        ActiveAdsCount.Set(float64(activeAdsCount))

        // Initialize AdvertisersCount
        var advertisersCount int64
        if err := database.DB.Model(&models.Advertiser{}).Count(&advertisersCount).Error; err != nil {
                log.Printf("Error counting advertisers: %v", err)
        }
        AdvertisersCount.Add(float64(advertisersCount)) // Counter metrics use Add instead of Set

        // Initialize PublishersCount
        var publishersCount int64
        if err := database.DB.Model(&models.Publisher{}).Count(&publishersCount).Error; err != nil {
                log.Printf("Error counting publishers: %v", err)
        }
        PublishersCount.Add(float64(publishersCount)) // Counter metrics use Add instead of Set

        // Initialize ImpressionCount
        var impressionCount int64
        if err := database.DB.Model(&models.AdsInteraction{}).Where("type = ?", int(models.Impression)).Count(&impressionCount).Error; err != nil {
                log.Printf("Error counting impressions: %v", err)
        }
        ImpressionCount.Add(float64(impressionCount))

        // Initialize ClickCount
        var clickCount int64
        if err := database.DB.Model(&models.AdsInteraction{}).Where("type = ?", int(models.Click)).Count(&clickCount).Error; err != nil {
                log.Printf("Error counting clicks: %v", err)
        }
        ClickCount.Add(float64(clickCount))

        // Initialize TotalRevenue
        var totalRevenue int64
        if err := database.DB.Model(&models.Transaction{}).Select("SUM(amount)").Where("income = ?", true).Scan(&totalRevenue).Error; err != nil {
                log.Printf("Error calculating total revenue: %v", err)
        }
        TotalRevenue.Add(float64(totalRevenue))

        // Initialize NumberOfBids and AverageBid
        var totalBids int64
        var sumBids int64
        if err := database.DB.Model(&models.Ad{}).Count(&totalBids).Error; err != nil {
                log.Printf("Error counting bids: %v", err)
        }
        if err := database.DB.Model(&models.Ad{}).Select("SUM(bid)").Scan(&sumBids).Error; err != nil {
                log.Printf("Error summing bids: %v", err)
        }
        NumberOfBids.Add(float64(totalBids))
        if totalBids > 0 {
                AverageBid.Set(float64(sumBids) / float64(totalBids))
        }

        // Initialize TotalAdvertiserBalance
        var totalAdvertiserBalance int64
        if err := database.DB.Model(&models.Advertiser{}).Select("SUM(balance)").Scan(&totalAdvertiserBalance).Error; err != nil {
                log.Printf("Error calculating total advertiser balance: %v", err)
        }
        TotalAdvertiserBalance.Set(float64(totalAdvertiserBalance))

        // Initialize TotalPublisherBalance
        var totalPublisherBalance int64
        if err := database.DB.Model(&models.Publisher{}).Select("SUM(balance)").Scan(&totalPublisherBalance).Error; err != nil {
                log.Printf("Error calculating total publisher balance: %v", err)
        }
        TotalPublisherBalance.Set(float64(totalPublisherBalance))

        // Initialize TransactionCount
        var successTransactions int64
        var failedTransactions int64
        if err := database.DB.Model(&models.Transaction{}).Where("successful = ?", true).Count(&successTransactions).Error; err != nil {
                log.Printf("Error counting successful transactions: %v", err)
        }
        if err := database.DB.Model(&models.Transaction{}).Where("successful = ?", false).Count(&failedTransactions).Error; err != nil {
                log.Printf("Error counting failed transactions: %v", err)
        }
        TransactionCount.WithLabelValues("success").Add(float64(successTransactions))
        TransactionCount.WithLabelValues("failure").Add(float64(failedTransactions))
}
