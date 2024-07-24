package handlers

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"encoding/base64"
	"net/http"

	"YellowBloomKnapsack/mini-yektanet/common/tokenhandler"
	"YellowBloomKnapsack/mini-yektanet/eventserver/worker"
	"YellowBloomKnapsack/mini-yektanet/eventserver/cache"
)

type EventServerHandler struct {
	tokenHandler  tokenhandler.TokenHandlerInterface
	workerService worker.WorkerInterface
	cacheService  cache.CacheInterface
}

func NewEventServerHandler(tokenHandler tokenhandler.TokenHandlerInterface, workerService worker.WorkerInterface, cacheService cache.CacheInterface) *EventServerHandler {
	workerService.Start()

	return &EventServerHandler{
		tokenHandler:  tokenHandler,
		workerService: workerService,
		cacheService:  cacheService,
	}
}

type TokenRequest struct {
	Token string `json:"token"`
}

func (h *EventServerHandler) PostClick(c *gin.Context) {
	privateKey := os.Getenv("PRIVATE_KEY")
	key, _ := base64.StdEncoding.DecodeString(privateKey)

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
		h.workerService.InvokeClickEvent(data, time.Now())
	}

	c.Redirect(http.StatusMovedPermanently, data.RedirectPath)
}

func (h *EventServerHandler) PostImpression(c *gin.Context) {
	privateKey := os.Getenv("PRIVATE_KEY")
	key, _ := base64.StdEncoding.DecodeString(privateKey)

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
		h.workerService.InvokeImpressionEvent(data, time.Now())
	}
}
