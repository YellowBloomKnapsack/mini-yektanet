package logic

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"YellowBloomKnapsack/mini-yektanet/common/dto"

	"github.com/stretchr/testify/assert"
)

func setupEnv() {
	os.Setenv("HOSTNAME", "localhost")
	os.Setenv("PANEL_PORT", "8080")
	os.Setenv("GET_ADS_API", "/ads")
	os.Setenv("ADS_FETCH_INTERVAL_SECS", "1")
}

func TestNewLogicService(t *testing.T) {
	setupEnv()

	interval, _ := strconv.Atoi(os.Getenv("ADS_FETCH_INTERVAL_SECS"))
	service := NewLogicService()

	ls, ok := service.(*LogicService)
	if !ok {
		t.Fatalf("expected LogicService, got %T", service)
	}

	assert.NotNil(t, ls)
	assert.Equal(t, 0, len(ls.adsList))
	assert.Equal(t, "http://localhost:8080/ads", ls.getAdsAPIPath)
	assert.Equal(t, interval, ls.interval)
}

func TestLogicService_GetBestAd_NoAds(t *testing.T) {
	setupEnv()

	service := NewLogicService()
	ls, _ := service.(*LogicService)

	ad, err := ls.GetBestAd()
	assert.Nil(t, ad)
	assert.NotNil(t, err)
	assert.Equal(t, "no ad was found", err.Error())
}

func TestLogicService_GetBestAd_WithAds(t *testing.T) {
	setupEnv()

	service := NewLogicService()
	ls, _ := service.(*LogicService)

	ads := []*dto.AdDTO{
		{ID: 1, Text: "Ad 1", Bid: 100},
		{ID: 2, Text: "Ad 2", Bid: 200},
		{ID: 3, Text: "Ad 3", Bid: 150},
	}

	ls.adsList = ads

	bestAd, err := ls.GetBestAd()
	assert.Nil(t, err)
	assert.NotNil(t, bestAd)
	assert.Equal(t, uint(2), bestAd.ID)
	assert.Equal(t, int64(200), bestAd.Bid)
}

func TestLogicService_UpdateAdsList(t *testing.T) {
	setupEnv()

	ads := []dto.AdDTO{
		{ID: 1, Text: "Ad 1", Bid: 100},
		{ID: 2, Text: "Ad 2", Bid: 200},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ads)
	}))
	defer server.Close()

	service := NewLogicService()
	ls, _ := service.(*LogicService)
	ls.getAdsAPIPath = server.URL

	err := ls.updateAdsList()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(ls.adsList))
	assert.Equal(t, uint(1), ls.adsList[0].ID)
	assert.Equal(t, uint(2), ls.adsList[1].ID)
}

func TestLogicService_StartTicker(t *testing.T) {
	setupEnv()

	ads := []dto.AdDTO{
		{ID: 1, Text: "Ad 1", Bid: 100},
		{ID: 2, Text: "Ad 2", Bid: 200},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ads)
	}))
	defer server.Close()

	service := NewLogicService()
	ls, _ := service.(*LogicService)
	ls.getAdsAPIPath = server.URL

	// Start the ticker
	ls.StartTicker()

	// Wait for the ticker to tick at least once
	time.Sleep(2 * time.Second)

	assert.Equal(t, 2, len(ls.adsList))
	assert.Equal(t, uint(1), ls.adsList[0].ID)
	assert.Equal(t, uint(2), ls.adsList[1].ID)
}
