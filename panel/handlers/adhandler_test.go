package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"

	"YellowBloomKnapsack/mini-yektanet/common/dto"
	"YellowBloomKnapsack/mini-yektanet/common/models"
	"YellowBloomKnapsack/mini-yektanet/panel/database"
)

const host string = "localhost:8090"

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
func TestGetActiveAdsSuccessGetAllAds(t *testing.T) {
	// setEnvVariables()

	// database.InitTestDB()

	r := gin.Default()
	r.GET("/", GetActiveAds)
	go r.Run(host) // the router needs to run in another goroutine
	time.Sleep(time.Millisecond * 100)
	resp, err := http.Get("http://" + host)

	require.NoError(t, err, "expected no error after making the request")

	// for testing purposes:
	// t.Log("\n\n\n\n\n\n\n")
	// t.Logf("%+v", resp)
	// t.Log("\n\n\n\n\n\n\n")
	// t.Logf("\n\n\n%T\n", resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error("could not read body of response")
	}

	var getActiveAdsAdDTOs []dto.AdDTO
	err = json.Unmarshal(body, &getActiveAdsAdDTOs)
	require.NoError(t, err, "json data sent from GetAllAds must be a slice of adDTO")

	var ads []models.Ad
	result := database.DB.Where("active = ?", true).Find(&ads)
	if result.Error != nil {
		t.Error("could not connect to database to check validity of data")
	}

	var adDTOs []dto.AdDTO
	for _, ad := range ads {
		adDTO := dto.AdDTO{
			ID:        ad.ID,
			Text:      ad.Text,
			ImagePath: "http://" + os.Getenv("PANEL_PUBLIC_HOSTNAME") + ":" + os.Getenv("PANEL_PORT") + ad.ImagePath,
			Bid:       ad.Bid,
			Website:   ad.Website,
		}
		adDTOs = append(adDTOs, adDTO)
	}
	if reflect.DeepEqual(adDTOs, getActiveAdsAdDTOs) == false {
		t.Error("returned ads by function getActiveAds is not the same as the locally read one")
	}
}

func setEnvVariables() {
	if err := godotenv.Load("../.env", "../../common/.env", "../../publisherwebsite/.env", "../../adserver/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
}
