package producer

import (
	"YellowBloomKnapsack/mini-yektanet/eventserver/grafana"
	"fmt"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ProducerInterface interface {
	Produce(protoData []byte, topicName string) error
}

type KafkaProducer struct {
	producer *kafka.Producer
}

func NewKafkaProducer() ProducerInterface {
	kafkaBootstrapServers := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	if kafkaBootstrapServers == "" {
		log.Fatal("KAFKA_BOOTSTRAP_SERVERS environment variable is not set")
	}

	// Initialize Kafka producer
	kafkaProducer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaBootstrapServers,
	})
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to create Kafka producer: %s", err))
	}

	return &KafkaProducer{
		producer: kafkaProducer,
	}
}

func (p *KafkaProducer) Produce(protoData []byte, topicName string) error {
	log.Println("Producing message to topic " + topicName)
	err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topicName, Partition: kafka.PartitionAny},
		Value:          protoData,
	}, nil)
	if err != nil {
		grafana.KafkaProducerMessages.WithLabelValues(topicName, "failure").Inc()
	} else {
		grafana.KafkaProducerMessages.WithLabelValues(topicName, "success").Inc()
	}
	p.producer.Flush(10000)
	return err
}
