package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"encoding/base64"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"YellowBloomKnapsack/mini-yektanet/common/dto"
)

// Custom Mock Implementations

type MockLogicService struct {
	BestAd *dto.AdDTO
	Err    error
}

func (m *MockLogicService) GetBestAd() (*dto.AdDTO, error) {
	return m.BestAd, m.Err
}

func (m *MockLogicService) StartTicker() {}

type MockTokenHandler struct {
	GenerateTokenResult string
	GenerateTokenError  error
}

func (m *MockTokenHandler) GenerateToken(interaction dto.InteractionType, adID uint, publisherUsername, redirectPath string, key []byte) (string, error) {
	return m.GenerateTokenResult, m.GenerateTokenError
}

func (m *MockTokenHandler) VerifyToken(encryptedToken string, key []byte) (*dto.CustomToken, error) {
	return nil, nil
}

func setupEnv() {
	os.Setenv("EVENT_SERVER_PORT", "8080")
	os.Setenv("HOSTNAME", "localhost")
	os.Setenv("CLICK_REQ_PATH", "click")
	os.Setenv("IMPRESSION_REQ_PATH", "impression")
	os.Setenv("PRIVATE_KEY", base64.StdEncoding.EncodeToString([]byte("mysecretkey123456")))
}

func TestGetAd(t *testing.T) {
	setupEnv()

	gin.SetMode(gin.TestMode)

	// Create mock instances
	mockLogicService := &MockLogicService{}
	mockTokenHandler := &MockTokenHandler{}

	// Create the handler with mocks
	handler := NewAdServerHandler(mockLogicService, mockTokenHandler)

	r := gin.Default()
	r.GET("/:publisherUsername", handler.GetAd)

	t.Run("success", func(t *testing.T) {
		ad := &dto.AdDTO{
			ID:        1,
			Text:      "Sample Ad",
			ImagePath: "http://example.com/image.png",
			Website:   "http://example.com",
			Bid:       100,
		}

		// Setup mocks for this test case
		mockLogicService.BestAd = ad
		mockTokenHandler.GenerateTokenResult = "token"

		req, _ := http.NewRequest("GET", "/testuser", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "http://example.com/image.png")
		assert.Contains(t, w.Body.String(), "token")
	})

	t.Run("no ad found", func(t *testing.T) {
		mockLogicService.BestAd = nil
		mockLogicService.Err = fmt.Errorf("no ad was found")

		req, _ := http.NewRequest("GET", "/testuser", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
