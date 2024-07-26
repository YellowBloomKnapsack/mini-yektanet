package logic

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"YellowBloomKnapsack/mini-yektanet/common/dto"

	"github.com/stretchr/testify/assert"
)

func setupEnv() {
	os.Setenv("PANEL_HOSTNAME", "localhost")
	os.Setenv("PANEL_PORT", "8080")
	os.Setenv("GET_ADS_API", "/ads")
	os.Setenv("ADS_FETCH_INTERVAL_SECS", "1")
}

func TestBestScoreOn(t *testing.T) {
	setupEnv()

	ls := NewLogicService().(*LogicService)

	ads := []*dto.AdDTO{
		{ID: 1, Score: 5},
		{ID: 2, Score: 10},
		{ID: 3, Score: 7},
	}

	bestAd := ls.bestScoreOn(ads)
	assert.Equal(t, uint(2), bestAd.ID)

	ads = []*dto.AdDTO{}

	bestAd = ls.bestScoreOn(ads)
	assert.Nil(t, bestAd)
}

func TestRandomOn(t *testing.T) {
	setupEnv()

	ls := NewLogicService().(*LogicService)

	ads := []*dto.AdDTO{
		{ID: 1},
		{ID: 2},
		{ID: 3},
	}

	randomAd := ls.randomOn(ads)
	assert.NotNil(t, randomAd)

	ads = []*dto.AdDTO{}

	randomAd = ls.bestScoreOn(ads)
	assert.Nil(t, randomAd)
}

func TestValidsOn(t *testing.T) {
	setupEnv()

	ls := NewLogicService().(*LogicService)
	ls.brakedAdIds[2] = struct{}{}
	ls.brakedAdIds[5] = struct{}{}
	ls.brakedAdIds[6] = struct{}{}

	ads := []*dto.AdDTO{
		{ID: 1},
		{ID: 2},
		{ID: 3},
		{ID: 4},
		{ID: 5},
		{ID: 6},
	}

	validAds := ls.validsOn(ads)
	assert.Equal(t, 3, len(validAds))
	assert.Equal(t, uint(1), validAds[0].ID)
	assert.Equal(t, uint(3), validAds[1].ID)
	assert.Equal(t, uint(4), validAds[2].ID)
}

func TestGetBestAd(t *testing.T) {
	setupEnv()
	os.Setenv("UNVISITED_CHANCE", "100")

	ls := NewLogicService().(*LogicService)

	ls.visitedAds = []*dto.AdDTO{
		{ID: 1, Score: 5},
	}

	ls.unvisitedAds = []*dto.AdDTO{
		{ID: 2, Score: 10},
	}

	bestAd, err := ls.GetBestAd()
	assert.NoError(t, err)
	assert.Equal(t, uint(2), bestAd.ID)

	os.Setenv("UNVISITED_CHANCE", "0")

	bestAd, err = ls.GetBestAd()
	assert.NoError(t, err)
	assert.Equal(t, uint(1), bestAd.ID)
}

func TestUpdateAdsList(t *testing.T) {
	setupEnv()

	ads := []dto.AdDTO{
		{ID: 1, Impressions: 0},
		{ID: 2, Impressions: 1},
	}
	body, _ := json.Marshal(ads)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(body)
	}))
	defer ts.Close()

	ls := NewLogicService().(*LogicService)
	ls.getAdsAPIPath = ts.URL

	err := ls.updateAdsList()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(ls.unvisitedAds))
	assert.Equal(t, 1, len(ls.visitedAds))
}

func TestStartTicker(t *testing.T) {
	setupEnv()

	ls := NewLogicService().(*LogicService)
	ls.interval = 1

	ads := []dto.AdDTO{
		{ID: 1, Impressions: 0},
		{ID: 2, Impressions: 1},
	}
	body, _ := json.Marshal(ads)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(body)
	}))
	defer ts.Close()

	ls.getAdsAPIPath = ts.URL

	ls.StartTicker()

	time.Sleep(2 * time.Second)

	assert.Equal(t, 1, len(ls.unvisitedAds))
	assert.Equal(t, 1, len(ls.visitedAds))
}

func TestBrakeAd(t *testing.T) {
	setupEnv()

	// Define test parameters
	adId := uint(1)
	duration := 200 * time.Millisecond

	service := NewLogicService()
	ls, _ := service.(*LogicService)

	ls.BrakeAd(adId, duration)

	// Check if adId is added to the map
	if _, found := ls.brakedAdIds[adId]; !found {
		t.Errorf("expected adId %d to be in the map", adId)
		return
	}

	// Wait for the specified duration plus a small buffer to ensure removal
	time.Sleep(duration + 100*time.Millisecond)

	// Check if adId is removed from the map
	if _, found := ls.brakedAdIds[adId]; found {
		t.Errorf("expected adId %d to be removed from the map", adId)
		return
	}
}
