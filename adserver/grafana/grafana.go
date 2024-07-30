package grafana

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	AdsTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ads_total",
		Help: "Total number of ads in the system",
	})

	AdsVisitedCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ads_visited_count",
		Help: "Number of ads in the visited category",
	})

	AdsUnvisitedCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ads_unvisited_count",
		Help: "Number of ads in the unvisited category",
	})

	AdsNewAddedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ads_new_added_total",
		Help: "Total number of new ads added to the system",
	})

	AdsFetchErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ads_fetch_errors_total",
		Help: "Total number of errors occurred while fetching ads",
	})

	AdSelectionMethodTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ad_selection_method_total",
			Help: "Count of how often each ad selection method is used",
		},
		[]string{"method"},
	)
	AdsBrakedCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ads_braked_count",
		Help: "Number of ads currently in the braked state",
	})

	AdRequestTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ad_request_total",
		Help: "Total number of ad requests received",
	})

	AdRequestDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "ad_request_duration_seconds",
		Help:    "Duration of ad request processing",
		Buckets: prometheus.DefBuckets,
	})

	AdServedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ad_served_total",
		Help: "Total number of ads successfully served",
	})

	AdNotFoundTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ad_not_found_total",
		Help: "Number of times no ad was found to serve",
	})

	TokenGenerationErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "token_generation_errors_total",
			Help: "Number of errors occurred during token generation",
		},
		[]string{"token_type"},
	)

	AdBrakeRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ad_brake_requests_total",
		Help: "Total number of ad brake requests received",
	})

	InvalidAdBrakeRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "invalid_ad_brake_requests_total",
		Help: "Number of invalid ad brake requests",
	})
)
