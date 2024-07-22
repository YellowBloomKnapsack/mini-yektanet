package handlers

import (
	"YellowBloomKnapsack/mini-yektanet/common/dto"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

var clickTokens = make(map[string]bool)
var impressionTokens = make(map[string]bool)

type EventPayload struct {
	AdID         string `json:"ad_id"`
	PublisherID  string `json:"publisher_id"`
	Token        string `json:"token"`
	RedirectPath string `json:"redirect_path"`
}

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
	data, err := verifyToken(token, key)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	_, present := clickTokens[token]
	if !present {
		clickTokens[token] = true

		url := "http://" + os.Getenv("HOSTNAME") + ":" + os.Getenv("PANEL_PORT") + os.Getenv("INTERACTION_CLICK_API")

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
	data, err := verifyToken(token, key)

	fmt.Println(data)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	_, present := impressionTokens[token]
	if !present {
		impressionTokens[token] = true
		// 1) request panel api
	}

	c.JSON(http.StatusOK, gin.H{"status": "impression registered"})
}

type InteractionType uint8

const (
	ImpressionType InteractionType = 0
	ClickType      InteractionType = 1
)

type CustomToken struct {
	Interaction       InteractionType `json:"interaction"`
	AdID              uint            `json:"ad_id"`
	PublisherUsername string          `json:"publisher_username"`
	RedirectPath      string          `json:"redirect_path"`
	CreatedAt         int64           `json:"created_at"`
}

func decrypt(encryptedData string, key []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func verifyToken(encryptedToken string, key []byte) (*CustomToken, error) {
	tokenBytes, err := decrypt(encryptedToken, key)
	if err != nil {
		return nil, err
	}

	var token CustomToken
	err = json.Unmarshal(tokenBytes, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
