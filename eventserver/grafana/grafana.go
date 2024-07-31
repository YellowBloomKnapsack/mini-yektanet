package grafana

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	TokenValidationTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "token_validation_total",
			Help: "Total number of token validations",
		},
		[]string{"status"},
	)

	KafkaProducerMessages = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_producer_messages_total",
			Help: "Total number of Kafka messages produced",
		},
		[]string{"topic", "status"},
	)

	CacheOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_operations_total",
			Help: "Total number of cache operations",
		},
		[]string{"operation", "status"},
	)
)
