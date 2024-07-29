package handlers

import (
	"YellowBloomKnapsack/mini-yektanet/common/models"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"YellowBloomKnapsack/mini-yektanet/common/dto"
)

func setupEnv() {
	os.Setenv("EVENT_SERVER_PORT", "8082")
	os.Setenv("REDIS_BF_INIT_SIZE", "6000000")
	os.Setenv("REDIS_BF_ERR_RATE", "0.01")
	os.Setenv("EVENT_SERVER_HOSTNAME", "localhost")
	os.Setenv("KAFKA_TOPIC_CLICK", "click_events")
	os.Setenv("KAFKA_TOPIC_IMPRESSION", "impression_events")
	os.Setenv("PRIVATE_KEY", "c2VjcmV0")
}

// Mock TokenHandler
type MockTokenHandler struct{}

func (m *MockTokenHandler) GenerateToken(interaction models.AdsInteractionType, adID, publisherID uint, bid int64, redirectPath string, key []byte) (string, error) {
	// Not needed for these tests
	return "duplicate", nil
}

func (m *MockTokenHandler) VerifyToken(encryptedToken string, key []byte) (*dto.CustomToken, error) {
	// Return a mock token
	return &dto.CustomToken{
		Interaction:  models.Click,
		AdID:         123,
		PublisherID:  234,
		RedirectPath: "http://example.com",
		CreatedAt:    time.Now().Unix(),
		Bid:          1000,
	}, nil
}

// Mock CacheService
type MockCacheService struct {
	times int
}

func (m *MockCacheService) IsPresent(token string) bool {
	return m.times > 0
}

func (m *MockCacheService) Add(token string) {
	m.times++
}

type MockProducerService struct {
	impCnt   int
	clickCnt int
}

func (p *MockProducerService) Produce(payload []byte, topic string) error {
	if topic == os.Getenv("KAFKA_TOPIC_CLICK") {
		p.clickCnt++
	} else if topic == os.Getenv("KAFKA_TOPIC_IMPRESSION") {
		p.impCnt++
	}
	return nil
}
func TestProduceImpressionIfTokenValid(t *testing.T) {
	setupEnv()

	mockTokenHandler := &MockTokenHandler{}
	mockCacheService := &MockCacheService{times: 0}
	mockProducerService := &MockProducerService{clickCnt: 0, impCnt: 0}

	handler := NewEventServerHandler(mockTokenHandler, mockCacheService, mockProducerService)

	data, _ := mockTokenHandler.VerifyToken("", []byte(""))

	handler.produceImpressionIfTokenValid("", data)

	assert.Equal(t, mockProducerService.clickCnt, 0)
	assert.Equal(t, mockProducerService.impCnt, 1)

	handler.produceImpressionIfTokenValid("", data)

	assert.Equal(t, mockProducerService.clickCnt, 0)
	assert.Equal(t, mockProducerService.impCnt, 1)
}

func TestProduceClickIfTokenValid(t *testing.T) {
	setupEnv()

	mockTokenHandler := &MockTokenHandler{}
	mockCacheService := &MockCacheService{times: 0}
	mockProducerService := &MockProducerService{clickCnt: 0, impCnt: 0}

	handler := NewEventServerHandler(mockTokenHandler, mockCacheService, mockProducerService)

	data, _ := mockTokenHandler.VerifyToken("", []byte(""))

	handler.produceClickIfTokenValid("", data)

	assert.Equal(t, mockProducerService.clickCnt, 1)
	assert.Equal(t, mockProducerService.impCnt, 0)

	handler.produceClickIfTokenValid("", data)

	assert.Equal(t, mockProducerService.clickCnt, 1)
	assert.Equal(t, mockProducerService.impCnt, 0)
}

func TestPostClick(t *testing.T) {
	setupEnv()

	mockTokenHandler := &MockTokenHandler{}
	mockCacheService := &MockCacheService{times: 0}
	mockProducerService := &MockProducerService{clickCnt: 0, impCnt: 0}

	handler := NewEventServerHandler(mockTokenHandler, mockCacheService, mockProducerService)

	r := gin.Default()
	r.POST("/click", handler.PostClick)

	req := httptest.NewRequest(http.MethodPost, "/click", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = ioutil.NopCloser(bytes.NewReader([]byte(`{"token":"dummy-token"}`)))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMovedPermanently, w.Code)
}

func TestPostImpression(t *testing.T) {
	setupEnv()

	mockTokenHandler := &MockTokenHandler{}
	mockCacheService := &MockCacheService{times: 0}
	mockProducerService := &MockProducerService{clickCnt: 0, impCnt: 0}

	handler := NewEventServerHandler(mockTokenHandler, mockCacheService, mockProducerService)

	r := gin.Default()
	r.POST("/impression", handler.PostImpression)

	req := httptest.NewRequest(http.MethodPost, "/impression", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = ioutil.NopCloser(bytes.NewReader([]byte(`{"token":"dummy-token"}`)))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
