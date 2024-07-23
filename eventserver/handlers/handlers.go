package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"YellowBloomKnapsack/mini-yektanet/common/dto"
	"YellowBloomKnapsack/mini-yektanet/common/tokeninterface"

	"github.com/gin-gonic/gin"
)

var clickTokens = make(map[string]bool)
var impressionTokens = make(map[string]bool)

type TokenRequest struct {
	Token string `json:"token"`
}

//func invokeClickEvent()

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

	//c.JSON(http.StatusOK, gin.H{"link": data.RedirectPath})
	c.Redirect(http.StatusMovedPermanently, data.RedirectPath)

	_, present := clickTokens[token]
	if !present {
		clickTokens[token] = true

		url := "http://" + os.Getenv("HOSTNAME") + ":" + os.Getenv("PANEL_PORT") + os.Getenv("INTERACTION_CLICK_API")

		go invokeClickEvent(data, url, time.Now())
	}

	// c.Redirect(http.StatusMovedPermanently, data.RedirectPath)
}

func invokeClickEvent(data *tokeninterface.CustomToken, url string, clickTime time.Time) {
	dataDto := dto.InteractionDto{
		PublisherUsername: data.PublisherUsername,
		ClickTime:         clickTime,
		AdID:              data.AdID,
	}

	jsonData, err := json.Marshal(dataDto)
	if err != nil {
		fmt.Printf("failed to marshal InteractionDto to JSON: %v\n", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("failed to send POST request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("failed to fetch ads: status code %d\n", resp.StatusCode)
	}
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
		clickTokens[token] = true

		url := "http://" + os.Getenv("HOSTNAME") + ":" + os.Getenv("PANEL_PORT") + os.Getenv("INTERACTION_IMPRESSION_API")

		dataDto := dto.InteractionDto{
			PublisherUsername: data.PublisherUsername,
			ClickTime:         time.Now(),
			AdID:              data.AdID,
		}

		jsonData, err := json.Marshal(dataDto)
		if err != nil {
			fmt.Printf("failed to marshal InteractionDto to JSON: %v\n", err)
			return
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("failed to send POST request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			fmt.Printf("failed to fetch ads: status code %d\n", resp.StatusCode)
		}
	}
}
