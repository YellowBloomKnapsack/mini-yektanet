package handlers

import (
	"encoding/base64"
	"net/http"
	"os"
	"time"

	"YellowBloomKnapsack/mini-yektanet/common/tokenhandler"
	"YellowBloomKnapsack/mini-yektanet/eventserver/worker"

	"github.com/gin-gonic/gin"
)

type EventServerHandler struct {
	clickTokens map[string]bool
	impressionTokens map[string]bool
	tokenHandler tokenhandler.TokenHandlerInterface
	workerService worker.WorkerInterface
}

func NewEventServerHandler(tokenHandler tokenhandler.TokenHandlerInterface, workerService worker.WorkerInterface) *EventServerHandler {
	workerService.Start()
	return &EventServerHandler{
		clickTokens:      make(map[string]bool),
		impressionTokens: make(map[string]bool),
		tokenHandler:     tokenHandler,
		workerService:    workerService,
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

	_, present := h.clickTokens[token]
	if !present {
		h.clickTokens[token] = true
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

	_, present := h.impressionTokens[token]
	if !present {
		h.impressionTokens[token] = true
		h.workerService.InvokeImpressionEvent(data, time.Now())
	}
}
