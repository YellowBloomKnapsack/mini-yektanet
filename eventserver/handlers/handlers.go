package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"

	"YellowBloomKnapsack/mini-yektanet/common/dto"
	"YellowBloomKnapsack/mini-yektanet/common/tokenhandler"
	"YellowBloomKnapsack/mini-yektanet/eventserver/worker"
	"YellowBloomKnapsack/mini-yektanet/common/cache"
)

type EventServerHandler struct {
	tokenHandler         tokenhandler.TokenHandlerInterface
	cacheService         cache.CacheInterface
	kafkaProducer        *kafka.Producer
	kafkaTopicClick      string
	kafkaTopicImpression string
}

// NewEventServerHandler initializes the event server handler with a Kafka producer.
func NewEventServerHandler(tokenHandler tokenhandler.TokenHandlerInterface, cacheService cache.CacheInterface) *EventServerHandler {
	kafkaBootstrapServers := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	if kafkaBootstrapServers == "" {
		panic("KAFKA_BOOTSTRAP_SERVERS environment variable is not set")
	}

	// Initialize Kafka producer
	kafkaProducer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaBootstrapServers,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to create Kafka producer: %s", err))
	}

	return &EventServerHandler{
		tokenHandler:         tokenHandler,
		cacheService:         cacheService,
		kafkaProducer:        kafkaProducer,
		kafkaTopicClick:      os.Getenv("KAFKA_TOPIC_CLICK"),
		kafkaTopicImpression: os.Getenv("KAFKA_TOPIC_IMPRESSION"),
	}
}

type TokenRequest struct {
	Token string `json:"token"`
}

// PostClick handles click events and produces them to a Kafka topic.
func (h *EventServerHandler) PostClick(c *gin.Context) {
	privateKey := os.Getenv("PRIVATE_KEY")
	key, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode private key"})
		return
	}

	var req TokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := req.Token
	data, err := h.tokenHandler.VerifyToken(token, key)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	present := h.cacheService.IsPresent(token)
	if !present {
		h.cacheService.Add(token)
		err = h.produceClickEvent(data, time.Now())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to produce click event"})
			return
		}
	}

	c.Redirect(http.StatusMovedPermanently, data.RedirectPath)
}

// PostImpression handles impression events and produces them to a Kafka topic.
func (h *EventServerHandler) PostImpression(c *gin.Context) {
	privateKey := os.Getenv("PRIVATE_KEY")
	key, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode private key"})
		return
	}

	var req TokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := req.Token
	data, err := h.tokenHandler.VerifyToken(token, key)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	present := h.cacheService.IsPresent(token)
	if !present {
		h.cacheService.Add(token)
		err = h.produceImpressionEvent(data, time.Now())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to produce impression event"})
			return
		}
	}
}

// produceClickEvent sends a click event to the Kafka topic.
func (h *EventServerHandler) produceClickEvent(data *dto.CustomToken, clickTime time.Time) error {
	eventData := &dto.ClickEvent{
		PublisherUsername: data.PublisherUsername,
		EventTime:         clickTime.Format(time.RFC3339),
		AdId:              uint32(data.AdID),
	}

	// Marshal to Protobuf
	protoData, err := proto.Marshal(eventData)
	if err != nil {
		fmt.Printf("Failed to marshal click event to protobuf: %v\n", err)
		return err
	}

	// Produce message to Kafka
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)

	err = h.kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &h.kafkaTopicClick, Partition: kafka.PartitionAny},
		Value:          protoData,
	}, deliveryChan)

	if err != nil {
		fmt.Printf("Failed to produce click event to Kafka: %v\n", err)
		return err
	}

	// Wait for message delivery
	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
		return m.TopicPartition.Error
	}

	fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
		*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	return nil
}

// produceImpressionEvent sends an impression event to the Kafka topic.
func (h *EventServerHandler) produceImpressionEvent(data *dto.CustomToken, impressionTime time.Time) error {
	eventData := &dto.ImpressionEvent{
		PublisherUsername: data.PublisherUsername,
		EventTime:         impressionTime.Format(time.RFC3339),
		AdId:              uint32(data.AdID),
	}

	// Marshal to Protobuf
	protoData, err := proto.Marshal(eventData)
	if err != nil {
		fmt.Printf("Failed to marshal impression event to protobuf: %v\n", err)
		return err
	}

	// Produce message to Kafka
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)

	err = h.kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &h.kafkaTopicImpression, Partition: kafka.PartitionAny},
		Value:          protoData,
	}, deliveryChan)

	if err != nil {
		fmt.Printf("Failed to produce impression event to Kafka: %v\n", err)
		return err
	}

	// Wait for message delivery
	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
		return m.TopicPartition.Error
	}

	fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
		*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	return nil
}
