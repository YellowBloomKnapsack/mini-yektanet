package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"YellowBloomKnapsack/mini-yektanet/common/models"
	"YellowBloomKnapsack/mini-yektanet/panel/database"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// os.Setenv("GOTEST_SEQUENTIAL", "1")
	cleanup := setupTest()
	exitCode := m.Run()
	cleanup()
	os.Exit(exitCode)
}

var sequentialTests = []func(t *testing.T){
	testAddFunds,
	testCreateAd,
	testHandleEditAd,
	testPublisherPanel,
	testWithdrawPublisherBalanceSuccessfulWithdrawal,
}

func TestSequentialTests(t *testing.T) {
	for _, fn := range sequentialTests {
		// name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
		// t.Run(name, fn)
		fn(t)
	}
}

// TestAdvertiserPanel tests the AdvertiserPanel handler
func TestAdvertiserPanel(t *testing.T) {
	// Initialize a new Gin router
	r := gin.Default()
	r.LoadHTMLGlob("../templates/*")
	r.GET("/advertiser/:username/panel", AdvertiserPanel)

	// Create a test advertiser
	advertiser := models.Advertiser{
		Username: "newtestadvertiser",
		Balance:  1000,
	}
	assert.NoError(t, database.DB.Create(&advertiser).Error)

	// Create a test request
	req, _ := http.NewRequest("GET", "/advertiser/newtestadvertiser/panel", nil)

	// Make the request and check the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "newtestadvertiser")
	// assert.Contains(t, w.Body.String(), fmt.Sprintf("%d", advertiser.Balance))
}

// TestAddFunds tests the AddFunds handler
func testAddFunds(t *testing.T) {
	r := gin.Default()
	r.POST("/advertiser/:username/add-funds", AddFunds)

	// Create a test advertiser
	advertiser := models.Advertiser{
		Username: "testadvertiser",
		Balance:  1000,
	}
	assert.NoError(t, database.DB.Create(&advertiser).Error)

	// Create a test request
	data := "amount=500"
	req, _ := http.NewRequest("POST", "/advertiser/testadvertiser/add-funds", bytes.NewBufferString(data))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Make the request and check the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusSeeOther, w.Code)

	// Verify if the advertiser's balance has been updated correctly
	var updatedAdvertiser models.Advertiser
	assert.NoError(t, database.DB.First(&updatedAdvertiser, advertiser.ID).Error)
	assert.Equal(t, int64(1500), updatedAdvertiser.Balance)

	// Verify if a new transaction has been created
	var transaction models.Transaction
	assert.NoError(t, database.DB.Where("customer_id = ? AND customer_type = ?", updatedAdvertiser.ID, models.Customer_Advertiser).First(&transaction).Error)
	assert.Equal(t, int64(500), transaction.Amount)
	assert.True(t, transaction.Successful)
}

// TestCreateAd tests the CreateAd handler
func testCreateAd(t *testing.T) {
	r := gin.Default()
	r.POST("/advertiser/:username/create-ad", CreateAd)

	// Create a test advertiser
	advertiser := models.Advertiser{
		Username: "testadvertiser4",
		Balance:  1000,
	}
	assert.NoError(t, database.DB.Create(&advertiser).Error)

	// Create a test file
	imagePath := "test_image.png"
	file, err := os.Create(imagePath)
	assert.NoError(t, err)
	file.Close()
	defer os.Remove(imagePath)

	// Prepare a form data request
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	writer.WriteField("title", "Test Ad")
	writer.WriteField("website", "https://example.com")
	writer.WriteField("bid", "100")

	part, err := writer.CreateFormFile("image", filepath.Base(imagePath))
	assert.NoError(t, err)
	testImageData := []byte("fake image data")
	part.Write(testImageData)

	writer.Close()

	// Create a test request
	req, _ := http.NewRequest("POST", "/advertiser/testadvertiser4/create-ad", &buffer)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Make the request and check the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusSeeOther, w.Code)

	// Verify if the ad has been created
	var ad models.Ad
	assert.NoError(t, database.DB.Where("advertiser_id = ? AND text = ?", advertiser.ID, "Test Ad").First(&ad).Error)
	assert.Equal(t, int64(100), ad.Bid)
	assert.Equal(t, "https://example.com", ad.Website)
}

// TestToggleAd tests the ToggleAd handler
func TestToggleAd(t *testing.T) {
	r := gin.Default()
	r.POST("/advertiser/:username/toggle-ad", ToggleAd)

	// Create a test advertiser and ad
	advertiser := models.Advertiser{
		Username: "testadvertiser5",
		Balance:  1000,
	}
	assert.NoError(t, database.DB.Create(&advertiser).Error)

	ad := models.Ad{
		Text:         "Test Ad",
		Bid:          100,
		AdvertiserID: advertiser.ID,
		Active:       true,
	}
	assert.NoError(t, database.DB.Create(&ad).Error)

	// Create a test request
	data := fmt.Sprintf("ad_id=%d", ad.ID)
	req, _ := http.NewRequest("POST", "/advertiser/testadvertiser5/toggle-ad", bytes.NewBufferString(data))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Make the request and check the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusSeeOther, w.Code)

	// Verify if the ad's active status has been toggled
	var updatedAd models.Ad
	assert.NoError(t, database.DB.First(&updatedAd, ad.ID).Error)
	assert.Equal(t, false, updatedAd.Active)

	// Toggle back
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusSeeOther, w.Code)

	assert.NoError(t, database.DB.First(&updatedAd, ad.ID).Error)
	assert.Equal(t, true, updatedAd.Active)
}

// TestAdReport tests the AdReport handler
func TestAdReport(t *testing.T) {
	r := gin.Default()
	r.LoadHTMLGlob("../templates/*")
	r.GET("/advertiser/ad-report/:id", AdReport)

	// Create a test advertiser and ad with interactions
	advertiser := models.Advertiser{
		Username: "testadvertiser6",
		Balance:  1000,
	}
	assert.NoError(t, database.DB.Create(&advertiser).Error)

	ad := models.Ad{
		Text:         "Test Ad",
		Bid:          100,
		AdvertiserID: advertiser.ID,
		Active:       true,
	}
	assert.NoError(t, database.DB.Create(&ad).Error)

	// Add some interactions
	interactions := []models.AdsInteraction{
		{Type: int(models.Impression), AdID: ad.ID},
		{Type: int(models.Impression), AdID: ad.ID},
		{Type: int(models.Click), AdID: ad.ID},
	}
	assert.NoError(t, database.DB.Create(&interactions).Error)

	// Create a test request
	req, _ := http.NewRequest("GET", fmt.Sprintf("/advertiser/ad-report/%d", ad.ID), nil)

	// Make the request and check the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "2")    // 2 Impressions
	assert.Contains(t, w.Body.String(), "1")    // 1 Click
	assert.Contains(t, w.Body.String(), "50.0") // CTR = 50.0
}

// TestHandleEditAd tests the HandleEditAd handler
func testHandleEditAd(t *testing.T) {
	r := gin.Default()
	r.POST("/advertiser/:username/edit-ad", HandleEditAd)

	// Create a test advertiser and ad
	advertiser := models.Advertiser{
		Username: "testadvertiser3",
		Balance:  1000,
	}
	assert.NoError(t, database.DB.Create(&advertiser).Error)

	ad := models.Ad{
		Text:         "Old Ad Text",
		Bid:          100,
		AdvertiserID: advertiser.ID,
		Active:       true,
		Website:      "https://oldsite.com",
	}
	assert.NoError(t, database.DB.Create(&ad).Error)

	// Prepare new image data
	newImagePath := "new_test_image.png"
	newFile, err := os.Create(newImagePath)
	assert.NoError(t, err)
	newFile.Close()
	defer os.Remove(newImagePath)

	// Prepare a form data request
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	writer.WriteField("ad_id", strconv.Itoa(int(ad.ID)))
	writer.WriteField("text", "Updated Ad Text")
	writer.WriteField("website", "https://newsite.com")
	writer.WriteField("bid", "150")

	part, err := writer.CreateFormFile("image", filepath.Base(newImagePath))
	assert.NoError(t, err)
	newImageData := []byte("new fake image data")
	part.Write(newImageData)

	writer.Close()

	// Create a test request
	req, _ := http.NewRequest("POST", "/advertiser/testadvertiser3/edit-ad", &buffer)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Make the request and check the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusSeeOther, w.Code)

	// Verify if the ad has been updated
	var updatedAd models.Ad
	assert.NoError(t, database.DB.First(&updatedAd, ad.ID).Error)
	assert.Equal(t, "Updated Ad Text", updatedAd.Text)
	assert.Equal(t, "https://newsite.com", updatedAd.Website)
	assert.Equal(t, int64(150), updatedAd.Bid)
}

// Utility function to prepare a JSON request body
func createJSONRequestBody(t *testing.T, data interface{}) io.Reader {
	body, err := json.Marshal(data)
	assert.NoError(t, err)
	return bytes.NewBuffer(body)
}
