package handlers

import (
	"encoding/base64"
	"net/http"
	"os"
	"time"

	"YellowBloomKnapsack/mini-yektanet/common/tokeninterface"
	"YellowBloomKnapsack/mini-yektanet/eventserver/eventschannel"

	"github.com/gin-gonic/gin"
)

var clickTokens = make(map[string]bool)
var impressionTokens = make(map[string]bool)

type TokenRequest struct {
	Token string `json:"token"`
}

func PostClick(c *gin.Context) {
	privateKey := os.Getenv("PRIVATE_KEY")
	key, _ := base64.StdEncoding.DecodeString(privateKey)

	var req TokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := req.Token
	data, err := tokeninterface.VerifyToken(token, key)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	_, present := clickTokens[token]
	if !present {
		clickTokens[token] = true
		eventschannel.InvokeClickEvent(data, time.Now())
	}

	c.Redirect(http.StatusMovedPermanently, data.RedirectPath)
}

func PostImpression(c *gin.Context) {
	privateKey := os.Getenv("PRIVATE_KEY")
	key, _ := base64.StdEncoding.DecodeString(privateKey)

	var req TokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := req.Token
	data, err := tokeninterface.VerifyToken(token, key)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	_, present := impressionTokens[token]
	if !present {
		impressionTokens[token] = true
		eventschannel.InvokeImpressionEvent(data, time.Now())
	}
}
