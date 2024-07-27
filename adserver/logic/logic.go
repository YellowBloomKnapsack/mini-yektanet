package logic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
	"time"

	"YellowBloomKnapsack/mini-yektanet/common/dto"
	"YellowBloomKnapsack/mini-yektanet/common/cache"
)

type LogicInterface interface {
	GetBestAd() (*dto.AdDTO, error)
	StartTicker()
	BrakeAd(adId uint)
}

type LogicService struct {
	visitedAds    []*dto.AdDTO
	unvisitedAds  []*dto.AdDTO
	brakedAdsCache cache.CacheInterface
	getAdsAPIPath string
	interval      int
}

func NewLogicService(cache cache.CacheInterface) LogicInterface {
	interval, _ := strconv.Atoi(os.Getenv("ADS_FETCH_INTERVAL_SECS"))
	return &LogicService{
		visitedAds:    make([]*dto.AdDTO, 0),
		unvisitedAds:  make([]*dto.AdDTO, 0),
		brakedAdsCache: cache,
		getAdsAPIPath: "http://" + os.Getenv("PANEL_HOSTNAME") + ":" + os.Getenv("PANEL_PORT") + os.Getenv("GET_ADS_API"),
		interval:      interval,
	}
}

func (ls *LogicService) isBetterThan(lhs, rhs *dto.AdDTO) bool {
	return rhs.Score > lhs.Score
}

func (ls *LogicService) bestScoreOn(ads []*dto.AdDTO) *dto.AdDTO {
	if len(ads) == 0 {
		return nil
	}

	bestAd := ads[0]
	for _, ad := range ads {
		if ls.isBetterThan(bestAd, ad) {
			bestAd = ad
		}
	}

	return bestAd
}

func (ls *LogicService) randomOn(ads []*dto.AdDTO) *dto.AdDTO {
	if len(ads) == 0 {
		return nil
	}

	return ads[rand.IntN(len(ads))]
}

func (ls *LogicService) isValid(ad *dto.AdDTO) bool {
	return !ls.brakedAdsCache.IsPresent(strconv.FormatUint(uint64(ad.ID), 10))
}

func (ls *LogicService) validsOn(ads []*dto.AdDTO) []*dto.AdDTO {
	result := make([]*dto.AdDTO, 0)
	for _, ad := range ads {
		if ls.isValid(ad) {
			result = append(result, ad)
		}
	}

	return result
}

func (ls *LogicService) GetBestAd() (*dto.AdDTO, error) {
	validVisitedsAds := ls.validsOn(ls.visitedAds)
	validUnvisitedsAds := ls.validsOn(ls.unvisitedAds)

	if len(validUnvisitedsAds) == 0 && len(validVisitedsAds) == 0 {
		return nil, fmt.Errorf("no ad was found")
	}

	if len(validUnvisitedsAds) == 0 {
		return ls.bestScoreOn(validVisitedsAds), nil
	}

	if len(validVisitedsAds) == 0 {
		return ls.randomOn(validUnvisitedsAds), nil
	}

	randomNumber := rand.Float32()
	unvisitedChance, _ := strconv.Atoi(os.Getenv("UNVISITED_CHANCE"))
	if randomNumber < float32(unvisitedChance)/100.0 {
		return ls.randomOn(validUnvisitedsAds), nil
	} else {
		return ls.bestScoreOn(validVisitedsAds), nil
	}
}

func (ls *LogicService) updateAdsList() error {
	fmt.Println("Fetching ads from panel API...")

	resp, err := http.Get(ls.getAdsAPIPath)
	if err != nil {
		return fmt.Errorf("failed to fetch ads: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch ads: status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	var ads []dto.AdDTO
	if err := json.Unmarshal(body, &ads); err != nil {
		return fmt.Errorf("failed to unmarshal ads: %v", err)
	}

	var newVisitedAds []*dto.AdDTO
	var newUnvisitedAds []*dto.AdDTO
	for _, ad := range ads {
		if ad.Impressions == 0 {
			newUnvisitedAds = append(newUnvisitedAds, &ad)
		} else {
			newVisitedAds = append(newVisitedAds, &ad)
		}
	}

	old_len := len(ls.unvisitedAds) + len(ls.visitedAds)
	new_len := len(newUnvisitedAds) + len(newVisitedAds)
	fmt.Printf("%d new ads were fetched.\n", new_len-old_len)

	ls.visitedAds = newVisitedAds
	ls.unvisitedAds = newUnvisitedAds

	return nil
}

func (ls *LogicService) StartTicker() {
	go func() {
		ticker := time.NewTicker(time.Duration(ls.interval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err := ls.updateAdsList()
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}()
}

func (ls *LogicService) BrakeAd(adId uint) {
	ls.brakedAdsCache.Add(strconv.FormatUint(uint64(adId), 10))
}
