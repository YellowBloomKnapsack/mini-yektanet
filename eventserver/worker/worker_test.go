package worker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"YellowBloomKnapsack/mini-yektanet/common/dto"
)

func setupTestEnvironment() func() {
	// Save original environment variables
	oldClickBuffSize := os.Getenv("CLICK_BUFF_SIZE")
	oldImpressionBuffSize := os.Getenv("IMPRESSION_BUFF_SIZE")
	oldHostName := os.Getenv("PANEL_HOSTNAME")
	oldPanelPort := os.Getenv("PANEL_PORT")
	oldClickApi := os.Getenv("INTERACTION_CLICK_API")
	oldImpressionApi := os.Getenv("INTERACTION_IMPRESSION_API")

	// Set test environment variables
	os.Setenv("CLICK_BUFF_SIZE", "10")
	os.Setenv("IMPRESSION_BUFF_SIZE", "10")
	os.Setenv("PANEL_HOSTNAME", "localhost")
	os.Setenv("PANEL_PORT", "8080")
	os.Setenv("INTERACTION_CLICK_API", "/click")
	os.Setenv("INTERACTION_IMPRESSION_API", "/impression")

	// Restore environment variables after test
	return func() {
		os.Setenv("CLICK_BUFF_SIZE", oldClickBuffSize)
		os.Setenv("IMPRESSION_BUFF_SIZE", oldImpressionBuffSize)
		os.Setenv("PANEL_HOSTNAME", oldHostName)
		os.Setenv("PANEL_PORT", oldPanelPort)
		os.Setenv("INTERACTION_CLICK_API", oldClickApi)
		os.Setenv("INTERACTION_IMPRESSION_API", oldImpressionApi)
	}
}

func TestWorkerService_InvokeClickEvent(t *testing.T) {
	teardown := setupTestEnvironment()
	defer teardown()

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/click" {
			t.Errorf("expected URL path '/click', got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		// Read and check the request body
		var body dto.InteractionDto
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}

		expected := dto.InteractionDto{
			PublisherUsername: "user1",
			EventTime:         time.Now(), // Time is not directly comparable
			AdID:              123,
		}
		if body != expected {
			t.Errorf("expected body %+v, got %+v", expected, body)
		}

		w.WriteHeader(http.StatusCreated)
	}))
	defer mockServer.Close()

	// Set the WorkerService API URLs to use the mock server
	worker := NewWorkerService().(*WorkerService)
	worker.clickApiUrl = mockServer.URL + "/click"
	worker.impressionApiUrl = mockServer.URL + "/impression"

	customToken := &dto.CustomToken{
		Interaction:       dto.ClickType,
		AdID:              123,
		PublisherUsername: "user1",
		RedirectPath:      "redirect/path",
		CreatedAt:         time.Now().Unix(),
	}
	clickTime := time.Now()

	// Invoke the click event
	worker.InvokeClickEvent(customToken, clickTime)

	// Allow some time for the goroutine to process
	time.Sleep(500 * time.Millisecond)
}

func TestWorkerService_InvokeImpressionEvent(t *testing.T) {
	teardown := setupTestEnvironment()
	defer teardown()

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/impression" {
			t.Errorf("expected URL path '/impression', got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		// Read and check the request body
		var body dto.InteractionDto
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}

		expected := dto.InteractionDto{
			PublisherUsername: "user1",
			EventTime:         time.Now(), // Time is not directly comparable
			AdID:              123,
		}
		if body != expected {
			t.Errorf("expected body %+v, got %+v", expected, body)
		}

		w.WriteHeader(http.StatusCreated)
	}))
	defer mockServer.Close()

	// Set the WorkerService API URLs to use the mock server
	worker := NewWorkerService().(*WorkerService)
	worker.clickApiUrl = mockServer.URL + "/click"
	worker.impressionApiUrl = mockServer.URL + "/impression"

	customToken := &dto.CustomToken{
		Interaction:       dto.ImpressionType,
		AdID:              123,
		PublisherUsername: "user1",
		RedirectPath:      "redirect/path",
		CreatedAt:         time.Now().Unix(),
	}
	impressionTime := time.Now()

	// Invoke the impression event
	worker.InvokeImpressionEvent(customToken, impressionTime)

	// Allow some time for the goroutine to process
	time.Sleep(500 * time.Millisecond)
}
