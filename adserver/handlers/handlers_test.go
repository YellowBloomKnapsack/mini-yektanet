package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

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

func (m *MockLogicService) BrakeAd(uint) {}

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
	os.Setenv("EVENT_SERVER_HOSTNAME", "localhost")
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

type MockLogicService_Brake struct {
	BrakeAdCalled   bool
	BrakeAdAdId     uint
	BrakeAdDuration time.Duration
}

func (m *MockLogicService_Brake) GetBestAd() (*dto.AdDTO, error) {
	return nil, nil
}

func (m *MockLogicService_Brake) StartTicker() {}

func (m *MockLogicService_Brake) BrakeAd(adId uint) {
	m.BrakeAdCalled = true
	m.BrakeAdAdId = adId
	brakeSeconds, _ := strconv.Atoi(os.Getenv("BRAKE_DURATION_SECS"))
	m.BrakeAdDuration = time.Duration(brakeSeconds)*time.Second
}

// TestBrakeAdHandler tests the BrakeAd handler function
func TestBrakeAdHandler(t *testing.T) {
	// Initialize Gin and create a test server
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Initialize the mock service
	mockLogicService := &MockLogicService_Brake{
		BrakeAdCalled: false,
	}

	handler := &AdServerHandler{
		logicService:  mockLogicService,
	}

	// Register the route and handler
	r.POST("/brakeAd/:adId", handler.BrakeAd)

	// Setup mock expectations
	adId := 123

	// Create a request to test the handler
	req, err := http.NewRequest("POST", "/brakeAd/"+strconv.Itoa(adId), nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	// Record the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

	// Assert that BrakeAd was called with the correct parameters
	assert.True(t, mockLogicService.BrakeAdCalled, "BrakeAd should have been called")
	assert.Equal(t, uint(adId), mockLogicService.BrakeAdAdId, "BrakeAd adId should match")
	// assert.Equal(t, handler.brakeDuration, mockLogicService.BrakeAdDuration, "BrakeAd duration should match")
}
