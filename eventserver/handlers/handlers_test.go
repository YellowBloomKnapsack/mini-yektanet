package handlers

import (
	"net/http"
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"YellowBloomKnapsack/mini-yektanet/common/dto"
)

// Mock TokenHandler
type MockTokenHandler struct{}

func (m *MockTokenHandler) GenerateToken(interaction dto.InteractionType, adID uint, publisherUsername, redirectPath string, key []byte) (string, error) {
	// Not needed for these tests
	return "duplicate", nil
}

func (m *MockTokenHandler) VerifyToken(encryptedToken string, key []byte) (*dto.CustomToken, error) {
	// Return a mock token
	return &dto.CustomToken{
		Interaction:       dto.ClickType,
		AdID:              123,
		PublisherUsername: "user1",
		RedirectPath:      "http://example.com",
		CreatedAt:         time.Now().Unix(),
	}, nil
}

// Mock WorkerService
type MockWorkerService struct {
	clickEventCalled      bool
	impressionEventCalled bool
	clickEventData        *dto.CustomToken
	impressionEventData   *dto.CustomToken
}

func (m *MockWorkerService) Start() {
	// Not needed for these tests
}

func (m *MockWorkerService) InvokeClickEvent(data *dto.CustomToken, clickTime time.Time) {
	m.clickEventCalled = true
	m.clickEventData = data
}

func (m *MockWorkerService) InvokeImpressionEvent(data *dto.CustomToken, impressionTime time.Time) {
	m.impressionEventCalled = true
	m.impressionEventData = data
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

func TestPostClick(t *testing.T) {
	privateKey := "c2VjcmV0" // base64 for 'secret'
	os.Setenv("PRIVATE_KEY", privateKey)

	mockTokenHandler := &MockTokenHandler{}
	mockWorkerService := &MockWorkerService{}
	mockCacheService := &MockCacheService{times: 0}

	handler := NewEventServerHandler(mockTokenHandler, mockWorkerService, mockCacheService)

	r := gin.Default()
	r.POST("/click", handler.PostClick)

	req := httptest.NewRequest(http.MethodPost, "/click", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = ioutil.NopCloser(bytes.NewReader([]byte(`{"token":"dummy-token"}`)))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMovedPermanently, w.Code)
	assert.True(t, mockWorkerService.clickEventCalled)
	assert.NotNil(t, mockWorkerService.clickEventData)
	assert.Equal(t, "user1", mockWorkerService.clickEventData.PublisherUsername)
	assert.Equal(t, uint(123), mockWorkerService.clickEventData.AdID)
	assert.Equal(t, "http://example.com", mockWorkerService.clickEventData.RedirectPath)

	// Request should not be send
	req = httptest.NewRequest(http.MethodPost, "/click", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = ioutil.NopCloser(bytes.NewReader([]byte(`{"token":"dummy-token"}`)))

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Body.Len(), 0)
}

func TestPostImpression(t *testing.T) {
	privateKey := "c2VjcmV0" // base64 for 'secret'
	os.Setenv("PRIVATE_KEY", privateKey)

	mockTokenHandler := &MockTokenHandler{}
	mockWorkerService := &MockWorkerService{}
	mockCacheService := &MockCacheService{times: 0}

	handler := NewEventServerHandler(mockTokenHandler, mockWorkerService, mockCacheService)

	r := gin.Default()
	r.POST("/impression", handler.PostImpression)

	req := httptest.NewRequest(http.MethodPost, "/impression", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = ioutil.NopCloser(bytes.NewReader([]byte(`{"token":"dummy-token"}`)))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, mockWorkerService.impressionEventCalled)
	assert.NotNil(t, mockWorkerService.impressionEventData)
	assert.Equal(t, "user1", mockWorkerService.impressionEventData.PublisherUsername)
	assert.Equal(t, uint(123), mockWorkerService.impressionEventData.AdID)

	// Request should not be send
	req = httptest.NewRequest(http.MethodPost, "/impression", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = ioutil.NopCloser(bytes.NewReader([]byte(`{"token":"dummy-token"}`)))

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Body.Len(), 0)
}

