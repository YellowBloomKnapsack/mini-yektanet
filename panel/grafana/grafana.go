package grafana

import (
        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
        ActiveAdsCount = promauto.NewGauge(prometheus.GaugeOpts{
                Name: "active_ads_count",
                Help: "The total number of active ads",
        })

        AdvertisersCount = promauto.NewGauge(prometheus.GaugeOpts{
                Name: "advertisers_count",
                Help: "The total number of advertisers",
        })

        PublishersCount = promauto.NewGauge(prometheus.GaugeOpts{
                Name: "publishers_count",
                Help: "The total number of publishers",
        })

        ImpressionCount = promauto.NewGauge(prometheus.GaugeOpts{
                Name: "impression_count",
                Help: "The total number of impressions",
        })

        ClickCount = promauto.NewGauge(prometheus.GaugeOpts{
                Name: "click_count",
                Help: "The total number of clicks",
        })

        TotalRevenue = promauto.NewGauge(prometheus.GaugeOpts{
                Name: "total_revenue",
                Help: "The total revenue generated",
        })

        NumberOfBids = promauto.NewGauge(prometheus.GaugeOpts{
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
        ActiveAdsCount.Set(0)
        AdvertisersCount.Set(0)
        PublishersCount.Set(0)
        ImpressionCount.Set(0)
        ClickCount.Set(0)
        TotalRevenue.Set(0)
        NumberOfBids.Set(0)
        AverageBid.Set(0)
        //TransactionCount // different setting needed for this
        TotalAdvertiserBalance.Set(0)
        TotalPublisherBalance.Set(0)
}