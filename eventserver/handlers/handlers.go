package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"

	"YellowBloomKnapsack/mini-yektanet/common/cache"
	"YellowBloomKnapsack/mini-yektanet/common/dto"
	"YellowBloomKnapsack/mini-yektanet/common/tokenhandler"
	"YellowBloomKnapsack/mini-yektanet/eventserver/producer"
)

type EventServerHandler struct {
	tokenHandler    tokenhandler.TokenHandlerInterface
	cacheService    cache.CacheInterface
	producer        producer.ProducerInterface
	clickTopic      string
	impressionTopic string
}

// NewEventServerHandler initializes the event server handler with a Kafka producer.
func NewEventServerHandler(tokenHandler tokenhandler.TokenHandlerInterface, cacheService cache.CacheInterface, producerService producer.ProducerInterface) *EventServerHandler {
	return &EventServerHandler{
		tokenHandler:    tokenHandler,
		cacheService:    cacheService,
		producer:        producerService,
		clickTopic:      os.Getenv("KAFKA_TOPIC_CLICK"),
		impressionTopic: os.Getenv("KAFKA_TOPIC_IMPRESSION"),
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

	fmt.Println("Token is hereeeeeeeeeee " + token)
	present := h.cacheService.IsPresent(token)
	fmt.Println(present)
	if !present {
		h.cacheService.Add(token)
		eventData := &dto.ClickEvent{
			PublisherId: uint32(data.PublisherID),
			EventTime:   time.Now().Format(time.RFC3339),
			AdId:        uint32(data.AdID),
			Bid:         data.Bid,
		}

		// Marshal to Protobuf
		protoData, err := proto.Marshal(eventData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to produce click event"})
		}

		err = h.producer.Produce(protoData, h.clickTopic)

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

		fmt.Println("meoooooooow and ")
		fmt.Println(data.PublisherID)
		eventData := &dto.ImpressionEvent{
			PublisherId: uint32(data.PublisherID),
			EventTime:   time.Now().Format(time.RFC3339),
			AdId:        uint32(data.AdID),
			Bid:         data.Bid,
		}

		// Marshal to Protobuf
		protoData, err := proto.Marshal(eventData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to produce click event"})
		}

		err = h.producer.Produce(protoData, h.impressionTopic)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to produce click event"})
			return
		}
	}
}
