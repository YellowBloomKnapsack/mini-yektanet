package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"YellowBloomKnapsack/mini-yektanet/common/dto"
	"YellowBloomKnapsack/mini-yektanet/common/models"
	"YellowBloomKnapsack/mini-yektanet/panel/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// Almost the same as the above, but this one is for single test instead of collection of tests
func setupTest() func() {
	godotenv.Load("../.env", "../../common/.env", "../../publisherwebsite/.env", "../../adserver/.env")
	os.Setenv("INTERACTION_CLICK_API", "/click")
	os.Setenv("INTERACTION_IMPRESSION_API", "/impression")
	os.Setenv("TEST_DB_PATH", "test.db")
	database.InitTestDB()
	fmt.Println("Test database initialized")
	return func() {
		os.Remove("test.db")
	}
}

func TestHandleClickAdInteraction(t *testing.T) {
	// Set up the test environment
	os.Setenv("YEKTANET_PORTION", "20")
	r := gin.Default()
	x := os.Getenv("INTERACTION_CLICK_API")
	r.POST(x, HandleClickAdInteraction)
	// Create a test publisher and ad
	publisher := models.Publisher{
		Username: "testpublisher",
		Balance:  1000,
	}
	assert.NoError(t, database.DB.Create(&publisher).Error)

	ad := models.Ad{
		Text:         "Test Ad",
		Bid:          100,
		AdvertiserID: 1,
		Advertiser: models.Advertiser{
			Balance: 1000,
		},
	}
	assert.NoError(t, database.DB.Create(&ad).Error)

	// Create a test request
	interactionDto := dto.InteractionDto{
		AdID:              ad.ID,
		PublisherUsername: publisher.Username,
		EventTime:         time.Now(),
	}
	body, _ := json.Marshal(interactionDto)
	req, _ := http.NewRequest("POST", os.Getenv("INTERACTION_CLICK_API"), bytes.NewBuffer(body))

	// Make the request and check the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Interaction recorded successfully", response["message"])

	// Check if the ad's total cost, publisher's balance, and advertiser's balance are updated correctly
	var updatedAd models.Ad
	assert.NoError(t, database.DB.First(&updatedAd, ad.ID).Error)
	assert.Equal(t, int64(100), updatedAd.TotalCost)

	var updatedPublisher models.Publisher
	assert.NoError(t, database.DB.First(&updatedPublisher, publisher.ID).Error)
	assert.Equal(t, int64(1080), updatedPublisher.Balance)

	var updatedAdvertiser models.Advertiser
	assert.NoError(t, database.DB.First(&updatedAdvertiser, ad.AdvertiserID).Error)
	assert.Equal(t, int64(900), updatedAdvertiser.Balance)
}

func testHandleImpressionAdInteraction(t *testing.T) {
	// Set up the test environment
	r := gin.Default()
	r.POST(os.Getenv("INTERACTION_IMPRESSION_API"), HandleImpressionAdInteraction)

	// Create a test publisher and ad
	publisher := models.Publisher{
		Username: "testpublisher2",
		Balance:  1000,
	}
	assert.NoError(t, database.DB.Create(&publisher).Error)

	ad := models.Ad{
		Text:         "Test Ad",
		Bid:          100,
		AdvertiserID: 1,
		Advertiser: models.Advertiser{
			Balance: 1000,
		},
	}
	assert.NoError(t, database.DB.Create(&ad).Error)

	// Create a test request
	interactionDto := dto.InteractionDto{
		AdID:              ad.ID,
		PublisherUsername: publisher.Username,
		EventTime:         time.Now(),
	}
	body, _ := json.Marshal(interactionDto)
	req, _ := http.NewRequest("POST", os.Getenv("INTERACTION_IMPRESSION_API"), bytes.NewBuffer(body))

	// Make the request and check the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Interaction recorded successfully", response["message"])

	// Check if the ad interaction is created correctly
	var createdInteraction models.AdsInteraction
	assert.NoError(t, database.DB.Where("ad_id = ? AND publisher_id = ?", ad.ID, publisher.ID).First(&createdInteraction).Error)
	assert.Equal(t, int(models.Impression), createdInteraction.Type)
}

func TestHandleClickAdInteractionWithInvalidRequest(t *testing.T) {
	// Set up the test environment
	os.Setenv("YEKTANET_PORTION", "20")
	r := gin.Default()
	x := os.Getenv("INTERACTION_CLICK_API")
	r.POST(x, HandleClickAdInteraction)

	// Create an invalid request
	body, _ := json.Marshal(map[string]interface{}{
		"ad_id":              "invalid_id",
		"publisher_username": "testpublisher20",
		"event_time":         time.Now(),
	})
	req, _ := http.NewRequest("POST", os.Getenv("INTERACTION_CLICK_API"), bytes.NewBuffer(body))

	// Make the request and check the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotEmpty(t, response["error"])
}

func TestHandleClickAdInteractionWithPublisherNotFound(t *testing.T) {
	// Set up the test environment
	os.Setenv("YEKTANET_PORTION", "20")
	r := gin.Default()
	x := os.Getenv("INTERACTION_CLICK_API")
	r.POST(x, HandleClickAdInteraction)

	// Create a test ad
	ad := models.Ad{
		Text:         "Test Ad",
		Bid:          100,
		AdvertiserID: 1,
		Advertiser: models.Advertiser{
			Balance: 1000,
		},
	}
	assert.NoError(t, database.DB.Create(&ad).Error)

	// Create a test request with a non-existent publisher
	interactionDto := dto.InteractionDto{
		AdID:              ad.ID,
		PublisherUsername: "non_existent_publisher",
		EventTime:         time.Now(),
	}
	body, _ := json.Marshal(interactionDto)
	req, _ := http.NewRequest("POST", os.Getenv("INTERACTION_CLICK_API"), bytes.NewBuffer(body))

	// Make the request and check the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Publisher not found", response["error"])
}
func TestHandleClickAdInteractionWithInvalidYektanetPortion(t *testing.T) {
	// Set up the test environment
	os.Setenv("YEKTANET_PORTION", "invalid_value")
	r := gin.Default()
	x := os.Getenv("INTERACTION_CLICK_API")
	r.POST(x, HandleClickAdInteraction)

	// Create a test publisher and ad
	publisher := models.Publisher{
		Username: "testpublisher21",
		Balance:  1000,
	}
	assert.NoError(t, database.DB.Create(&publisher).Error)

	ad := models.Ad{
		Text:         "Test Ad",
		Bid:          100,
		AdvertiserID: 1,
		Advertiser: models.Advertiser{
			Balance: 1000,
		},
	}
	assert.NoError(t, database.DB.Create(&ad).Error)

	// Create a test request
	interactionDto := dto.InteractionDto{
		AdID:              ad.ID,
		PublisherUsername: publisher.Username,
		EventTime:         time.Now(),
	}
	body, _ := json.Marshal(interactionDto)
	req, _ := http.NewRequest("POST", os.Getenv("INTERACTION_CLICK_API"), bytes.NewBuffer(body))

	// Make the request and check the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "Error parsing YEKTANET_PORTION environment variable")
}
