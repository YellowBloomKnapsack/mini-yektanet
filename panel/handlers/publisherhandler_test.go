package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"YellowBloomKnapsack/mini-yektanet/common/models"
	"YellowBloomKnapsack/mini-yektanet/panel/database"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPublisherPanel(t *testing.T) {
	r := gin.Default()
	r.LoadHTMLGlob("../templates/*")
	r.GET("/publisher/:username/panel", PublisherPanel)

	// Test Case 1: Existing Publisher
	t.Run("Existing Publisher", func(t *testing.T) {
		// Create a test publisher
		publisher := models.Publisher{
			Username: "testpublisher",
			Balance:  500,
		}
		assert.NoError(t, database.DB.Create(&publisher).Error)

		// Create a test ad interaction
		interaction := models.AdsInteraction{
			Type:        int(models.Click),
			EventTime:   time.Now(),
			Bid:         100,
			AdID:        1,
			PublisherID: publisher.ID,
		}
		assert.NoError(t, database.DB.Create(&interaction).Error)

		req, _ := http.NewRequest("GET", "/publisher/testpublisher/panel", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "testpublisher")
		assert.Contains(t, w.Body.String(), "500")
	})

	// Test Case 2: New Publisher
	t.Run("New Publisher", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/publisher/newpublisher/panel", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var newPublisher models.Publisher
		assert.NoError(t, database.DB.Where("username = ?", "newpublisher").First(&newPublisher).Error)
		assert.Equal(t, "newpublisher", newPublisher.Username)
		assert.Equal(t, int64(0), newPublisher.Balance)
	})
}
func TestWithdrawPublisherBalanceSuccessfulWithdrawal(t *testing.T) {
	r := gin.Default()
	r.POST("/publisher/:username/withdraw", WithdrawPublisherBalance)

	// Test Case 1: Successful Withdrawal
	publisher := models.Publisher{
		Username: "testpublisher",
		Balance:  1000,
	}
	assert.NoError(t, database.DB.Create(&publisher).Error)

	formData := url.Values{}
	formData.Set("amount", "200")

	req, _ := http.NewRequest("POST", "/publisher/testpublisher/withdraw", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, fmt.Sprintf("Withdrawn amount: %d", 200), response["message"])
	assert.Equal(t, float64(800), response["newBalance"])

	var updatedPublisher models.Publisher
	err := database.DB.Where("username = ?", "testpublisher").First(&updatedPublisher).Error
	assert.NoError(t, err)
	assert.Equal(t, int64(800), updatedPublisher.Balance)
}

func TestWithdrawPublisherBalance(t *testing.T) {
	r := gin.Default()
	r.POST("/publisher/:username/withdraw", WithdrawPublisherBalance)

	// Test Case 2: Invalid Withdrawal Amount (e.g., greater than balance)
	t.Run("Invalid Withdrawal Amount - Greater Than Balance", func(t *testing.T) {
		publisher := models.Publisher{
			Username: "richpublisher",
			Balance:  1000,
		}
		assert.NoError(t, database.DB.Create(&publisher).Error)

		formData := url.Values{}
		formData.Set("amount", "1500")

		req, _ := http.NewRequest("POST", "/publisher/richpublisher/withdraw", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Invalid withdrawal amount", response["error"])

		var updatedPublisher models.Publisher
		assert.NoError(t, database.DB.First(&updatedPublisher, publisher.ID).Error)
		assert.Equal(t, int64(1000), updatedPublisher.Balance) // Balance should remain the same
	})

	// Test Case 3: Invalid Withdrawal Amount (e.g., non-positive)
	t.Run("Invalid Withdrawal Amount - Non-Positive", func(t *testing.T) {
		publisher := models.Publisher{
			Username: "averagepublisher",
			Balance:  1000,
		}
		assert.NoError(t, database.DB.Create(&publisher).Error)

		formData := url.Values{}
		formData.Set("amount", "-100")

		req, _ := http.NewRequest("POST", "/publisher/averagepublisher/withdraw", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Invalid withdrawal amount", response["error"])

		var updatedPublisher models.Publisher
		assert.NoError(t, database.DB.First(&updatedPublisher, publisher.ID).Error)
		assert.Equal(t, int64(1000), updatedPublisher.Balance) // Balance should remain the same
	})

	// Test Case 4: Non-existent Publisher
	t.Run("Non-existent Publisher", func(t *testing.T) {
		formData := map[string]string{
			"amount": "100",
		}
		body, _ := json.Marshal(formData)

		req, _ := http.NewRequest("POST", "/publisher/nonexistent/withdraw", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Publisher not found", response["error"])
	})
}