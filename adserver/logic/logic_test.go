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

// Mock KVStorage
type MockKVStorage struct {
	kv map[string]string
}

func (s *MockKVStorage) Get(key string) (string, error) {
	return s.kv[key], nil
}

func (s *MockKVStorage) Set(key, value string) error {
	s.kv[key] = value
	return nil
}

// Mock CacheService
type MockCacheService struct {
	mark map[string]interface{}
}

func (m *MockCacheService) IsPresent(key string) bool {
	_, ok := m.mark[key]
	return ok
}

func (m *MockCacheService) Add(key string) {
	m.mark[key] = ""
	go func() {
		brakeSeconds, _ := strconv.Atoi(os.Getenv("BRAKE_DURATION_SECS"))
		time.Sleep(time.Duration(brakeSeconds))
		delete(m.mark, key)
	}()
}

func setupEnv() {
	os.Setenv("PANEL_HOSTNAME", "localhost")
	os.Setenv("PANEL_PORT", "8080")
	os.Setenv("GET_ADS_API", "/ads")
	os.Setenv("ADS_FETCH_INTERVAL_SECS", "1")
	os.Setenv("BRAKE_DURATION_SECS", "5")
}

func TestBestScoreOn(t *testing.T) {
	setupEnv()

	cache := &MockCacheService{make(map[string]interface{})}
	kv := &MockKVStorage{make(map[string]string)}
	ls := NewLogicService(cache, kv).(*LogicService)

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

	cache := &MockCacheService{make(map[string]interface{})}
	kv := &MockKVStorage{make(map[string]string)}
	ls := NewLogicService(cache, kv).(*LogicService)

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

	cache := &MockCacheService{make(map[string]interface{})}
	kv := &MockKVStorage{make(map[string]string)}
	ls := NewLogicService(cache, kv).(*LogicService)
	ls.brakedAdsCache.Add("2")
	ls.brakedAdsCache.Add("5")
	ls.brakedAdsCache.Add("6")

	kv.Set("1", "kw2")

	ads := []*dto.AdDTO{
		{ID: 1, Keywords: []string{"kw1", "kw2"}},
		{ID: 2, Keywords: []string{}},
		{ID: 3, Keywords: []string{"kw2", "kw3"}},
		{ID: 4, Keywords: []string{"kw4", "kw5"}},
		{ID: 5, Keywords: []string{}},
		{ID: 6, Keywords: []string{}},
	}

	validAds := ls.validsOn(ads, uint(1))
	assert.Equal(t, 2, len(validAds))
	assert.Equal(t, uint(1), validAds[0].ID)
	assert.Equal(t, uint(3), validAds[1].ID)
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

	cache := &MockCacheService{make(map[string]interface{})}
	kv := &MockKVStorage{make(map[string]string)}
	ls := NewLogicService(cache, kv).(*LogicService)
	ls.getAdsAPIPath = ts.URL

	err := ls.updateAdsList()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(ls.unvisitedAds))
	assert.Equal(t, 1, len(ls.visitedAds))
}

func TestStartTicker(t *testing.T) {
	setupEnv()

	cache := &MockCacheService{make(map[string]interface{})}
	kv := &MockKVStorage{make(map[string]string)}
	ls := NewLogicService(cache, kv).(*LogicService)
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

	cache := &MockCacheService{make(map[string]interface{})}
	kv := &MockKVStorage{make(map[string]string)}
	ls := NewLogicService(cache, kv).(*LogicService)

	ls.BrakeAd(adId)

	// Check if adId is added to the map
	if !ls.brakedAdsCache.IsPresent(strconv.FormatUint(uint64(adId), 10)) {
		t.Errorf("expected adId %d to be in the map", adId)
		return
	}

	// Wait for the specified duration plus a small buffer to ensure removal
	time.Sleep(duration + 100*time.Millisecond)

	// Check if adId is removed from the map
	if ls.brakedAdsCache.IsPresent(strconv.FormatUint(uint64(adId), 10)) {
		t.Errorf("expected adId %d to be removed from the map", adId)
		return
	}
}

func TestGetBestAd_NoAdsAvailable(t *testing.T) {
	cache := &MockCacheService{make(map[string]interface{})}
	kv := &MockKVStorage{make(map[string]string)}
	ls := NewLogicService(cache, kv).(*LogicService)

	_, err := ls.GetBestAd(uint(1))
	assert.Error(t, err)
	assert.Equal(t, "no ad was found", err.Error())
}

func TestGetBestAd_OnlyUnvisitedAdsAvailable(t *testing.T) {
	os.Setenv("UNVISITED_CHANCE", "100")

	cache := &MockCacheService{make(map[string]interface{})}
	kv := &MockKVStorage{make(map[string]string)}
	ls := NewLogicService(cache, kv).(*LogicService)
	ls.unvisitedAds = []*dto.AdDTO{
		{ID: 1, Score: 5},
		{ID: 2, Score: 10},
	}

	bestAd, err := ls.GetBestAd(uint(1))
	assert.NoError(t, err)
	assert.Contains(t, []uint{1, 2}, bestAd.ID)
}

func TestGetBestAd_OnlyVisitedAdsAvailable(t *testing.T) {
	cache := &MockCacheService{make(map[string]interface{})}
	kv := &MockKVStorage{make(map[string]string)}
	ls := NewLogicService(cache, kv).(*LogicService)
	ls.visitedAds = []*dto.AdDTO{
		{ID: 1, Score: 5},
		{ID: 2, Score: 10},
	}

	bestAd, err := ls.GetBestAd(uint(1))
	assert.NoError(t, err)
	assert.Equal(t, uint(2), bestAd.ID)
}

func TestGetBestAd_BothAdsAvailable_UnvisitedChance100(t *testing.T) {
	os.Setenv("UNVISITED_CHANCE", "100")

	cache := &MockCacheService{make(map[string]interface{})}
	kv := &MockKVStorage{make(map[string]string)}
	ls := NewLogicService(cache, kv).(*LogicService)
	ls.visitedAds = []*dto.AdDTO{
		{ID: 1, Score: 5},
	}
	ls.unvisitedAds = []*dto.AdDTO{
		{ID: 2, Score: 10},
	}

	bestAd, err := ls.GetBestAd(uint(1))
	assert.NoError(t, err)
	assert.Equal(t, uint(2), bestAd.ID)
}

func TestGetBestAd_BothAdsAvailable_UnvisitedChance0(t *testing.T) {
	os.Setenv("UNVISITED_CHANCE", "0")

	cache := &MockCacheService{make(map[string]interface{})}
	kv := &MockKVStorage{make(map[string]string)}
	ls := NewLogicService(cache, kv).(*LogicService)
	ls.visitedAds = []*dto.AdDTO{
		{ID: 1, Score: 5},
	}
	ls.unvisitedAds = []*dto.AdDTO{
		{ID: 2, Score: 10},
	}

	bestAd, err := ls.GetBestAd(uint(1))
	assert.NoError(t, err)
	assert.Equal(t, uint(1), bestAd.ID)
}

func TestGetBestAd_ValidUnvisitedAdsAvailable(t *testing.T) {
	os.Setenv("UNVISITED_CHANCE", "100")

	cache := &MockCacheService{make(map[string]interface{})}
	kv := &MockKVStorage{make(map[string]string)}
	ls := NewLogicService(cache, kv).(*LogicService)
	ls.unvisitedAds = []*dto.AdDTO{
		{ID: 1, Score: 5},
		{ID: 2, Score: 10},
	}

	// Braking one of the ads to make it invalid
	ls.BrakeAd(1)

	bestAd, err := ls.GetBestAd(uint(1))
	assert.NoError(t, err)
	assert.Equal(t, uint(2), bestAd.ID)
}

func TestGetBestAd_ValidVisitedAdsAvailable(t *testing.T) {
	cache := &MockCacheService{make(map[string]interface{})}
	kv := &MockKVStorage{make(map[string]string)}
	ls := NewLogicService(cache, kv).(*LogicService)
	ls.visitedAds = []*dto.AdDTO{
		{ID: 1, Score: 5},
		{ID: 2, Score: 10},
	}

	// Braking one of the ads to make it invalid
	ls.BrakeAd(2)

	bestAd, err := ls.GetBestAd(uint(1))
	assert.NoError(t, err)
	assert.Equal(t, uint(1), bestAd.ID)
}

func TestGetBestAd_ValidUnvisitedAndVisitedAdsAvailable(t *testing.T) {
	os.Setenv("UNVISITED_CHANCE", "50")

	cache := &MockCacheService{make(map[string]interface{})}
	kv := &MockKVStorage{make(map[string]string)}
	ls := NewLogicService(cache, kv).(*LogicService)
	ls.visitedAds = []*dto.AdDTO{
		{ID: 1, Score: 5},
	}
	ls.unvisitedAds = []*dto.AdDTO{
		{ID: 2, Score: 10},
	}

	// Braking the unvisited ad to make it invalid
	ls.BrakeAd(2)

	bestAd, err := ls.GetBestAd(uint(1))
	assert.NoError(t, err)
	assert.Equal(t, uint(1), bestAd.ID)
}

func TestHasIntersection(t *testing.T) {
	ls := &LogicService{}

	tests := []struct {
		lhs      []string
		rhs      []string
		expected bool
	}{
		{[]string{"apple", "banana", "cherry"}, []string{"grape", "banana", "kiwi"}, true},
		{[]string{"apple", "banana", "cherry"}, []string{"grape", "kiwi"}, false},
		{[]string{"apple", "banana", "cherry"}, []string{"apple", "kiwi"}, true},
		{[]string{}, []string{"grape", "kiwi"}, false},
		{[]string{"apple", "banana"}, []string{}, false},
		{[]string{}, []string{}, false},
		{[]string{"apple", "banana", "banana"}, []string{"banana"}, true},
		{[]string{"apple", "banana", "cherry"}, []string{"cherry", "banana"}, true},
	}

	for _, test := range tests {
		result := ls.hasIntersection(test.lhs, test.rhs)
		if result != test.expected {
			t.Errorf("hasIntersection(%v, %v) = %v; expected %v", test.lhs, test.rhs, result, test.expected)
		}
	}
}
