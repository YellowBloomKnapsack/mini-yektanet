package handlers

import (
	"os"
	"testing"
	"log"
	"reflect"
	"io/ioutil"
	"net/http"
	"encoding/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"YellowBloomKnapsack/mini-yektanet/common/dto"
	"YellowBloomKnapsack/mini-yektanet/common/models"
	"YellowBloomKnapsack/mini-yektanet/panel/database"
)

const host string = "localhost:8090"

func TestGetActiveAdsErrorNoDB(t *testing.T) {
	setEnvVariables()
	
	database.InitTestDB() // init the database (so database.DB - being used in adHandler - is not nil)
	db, _ := database.DB.DB()
	db.Close() // close database connection for the sake of testing	

	r := gin.Default()
	r.GET("/", GetActiveAds)
	go r.Run(host) // the router needs to run in another goroutine

	resp, err := http.Get("http://"+host)

	// for testing purposes:
	// t.Log("\n\n\n\n\n\n\n")
	// t.Logf("%+v", resp)
	// t.Log("\n\n\n\n\n\n\n")

	require.NoError(t, err, "expected no error after making the request")
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "response status code must be 500 (internal server error)")
}

func TestGetActiveAdsSuccessGetAllAds(t *testing.T) {
	setEnvVariables()

	database.InitTestDB()

	r := gin.Default()	
	r.GET("/", GetActiveAds)
	go r.Run(host) // the router needs to run in another goroutine	

	resp, err := http.Get("http://"+host)

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
			ImagePath: "http://" + os.Getenv("HOSTNAME") + ":" + os.Getenv("PANEL_PORT") + ad.ImagePath,
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