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

// "math/rand"
// "net/http"
// "time"

// "YellowBloomKnapsack/mini-yektanet/common/database"
// "YellowBloomKnapsack/mini-yektanet/common/models"

var clickTokens = make(map[string]bool, 0)
var impressionTokens = make(map[string]bool, 0)

type EventPayload struct {
	AdID         string `json:"ad_id"`
	PublisherID  string `json:"publisher_id"`
	Token        string `json:"token"`
	RedirectPath string `json:"redirect_path"`
}

//
//func convertCustomTokenToInteractionDto(token CustomToken) InteractionDto {
//	return
//}

func PostClick(c *gin.Context) {

	privateKey := os.Getenv("PRIVATE_KEY")
	key, _ := base64.StdEncoding.DecodeString(privateKey)

	token := c.Query("token")
	data, err := verifyToken(token, key)
	if err != nil {
		return
	}

	fmt.Println(data)

	_, present := clickTokens[token]
	if !present {
		clickTokens[token] = true
		fmt.Println("jadid")

		url := "http://" + os.Getenv("HOSTNAME") + ":" + os.Getenv("PANEL_PORT") + os.Getenv("INTERACTION_CLICK_API")

		fmt.Println(url)

		dataDto := dto.InteractionDto{
			PublisherUsername: data.PublisherUsername,
			ClickTime:         time.Now(),
			AdID:              data.AdID,
		}

		jsonData, err := json.Marshal(dataDto)

		if err != nil {
			fmt.Errorf("failed to marshal InteractionDto to JSON: %v", err)
			return
		}

		// Send HTTP POST request with JSON data
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Errorf("failed to send POST request: %v", err)
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			fmt.Errorf("failed to fetch ads: status code %d", resp.StatusCode)
		}

	}

	// redirect anyways
	c.Redirect(http.StatusMovedPermanently, data.RedirectPath)
}

func PostImpression(c *gin.Context) {
	privateKey := os.Getenv("PRIVATE_KEY")
	key, _ := base64.StdEncoding.DecodeString(privateKey)

	token := c.Query("token")
	data, err := verifyToken(token, key)
	if err != nil {
		return
	}

	fmt.Println(data)

	_, present := clickTokens[token]
	if !present {
		clickTokens[token] = true
		fmt.Println("jadid")
		// 1) request panel api
	}

	//c.JSON(http.StatusOK, gin.H{"status": "impression registered"})

}

type InteractionType uint8

const (
	ImpressionType InteractionType = 0
	ClickType      InteractionType = 1
)

// CustomToken defines the structure of the token
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

	//maybe add later
	//currentTime := time.Now().Unix()
	//if currentTime < token.CreatedAt {
	//	return nil, errors.New("token is not valid at this time")
	//}

	return &token, nil
}

//
//func main() {
//	r := gin.Default()
//
//	publicKey := os.Getenv("PUBLIC_KEY")
//	if publicKey == "" {
//		log.Fatal("PUBLIC_KEY environment variable is required")
//	}
//
//	key, err := base64.StdEncoding.DecodeString(publicKey)
//	if err != nil {
//		log.Fatal("Invalid PUBLIC_KEY format: must be base64 encoded")
//	}
//
//	r.GET("/verify_token", func(c *gin.Context) {
//		token := c.Query("token")
//
//		claims, err := verifyToken(token, key)
//		if err != nil {
//			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
//			return
//		}
//
//		c.JSON(http.StatusOK, gin.H{
//			"ad_id":         claims.AdID,
//			"publisher_id":  claims.PublisherUsername,
//			"redirect_path": claims.RedirectPath,
//			"created_at":    claims.CreatedAt,
//		})
//	})
//
//	r.Run(":8081")
//}
