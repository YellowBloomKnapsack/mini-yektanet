package handlers

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"

	"YellowBloomKnapsack/mini-yektanet/eventserver/grafana"
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
		grafana.TokenValidationTotal.WithLabelValues("invalid").Inc()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}
	grafana.TokenValidationTotal.WithLabelValues("valid").Inc()

	// Running in goroutine so the server wouldn't have to wait
	go h.produceClickIfTokenValid(token, data)

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

	// Running in goroutine so the server wouldn't have to wait
	go h.produceImpressionIfTokenValid(token, data)
}

func (h *EventServerHandler) produceImpressionIfTokenValid(token string, data *dto.CustomToken) {
	present := h.cacheService.IsPresent(token)
	if present {
		log.Printf("Token %s already present", token)
		return
	}

	h.cacheService.Add(token)

	eventData := &dto.ImpressionEvent{
		PublisherId: uint32(data.PublisherID),
		EventTime:   time.Now().Format(time.RFC3339),
		AdId:        uint32(data.AdID),
		Bid:         data.Bid,
	}

	// Marshal to Protobuf
	protoData, err := proto.Marshal(eventData)
	if err != nil {
		log.Println("Failed to marshal event data:", err)
		return
	}

	// Produce event
	err = h.producer.Produce(protoData, h.impressionTopic)
	if err != nil {
		log.Println("Failed to produce impression:", err)
		return
	}
}

func (h *EventServerHandler) produceClickIfTokenValid(token string, data *dto.CustomToken) {
	present := h.cacheService.IsPresent(token)
	if present {
		log.Printf("Token %s already present", token)
		return
	}

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
		log.Println("Failed to marshal event data:", err)
		return
	}

	err = h.producer.Produce(protoData, h.clickTopic)
	if err != nil {
		log.Println("Failed to produce event:", err)
		return
	}
}
