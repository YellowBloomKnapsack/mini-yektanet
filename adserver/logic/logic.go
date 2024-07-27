package logic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"YellowBloomKnapsack/mini-yektanet/common/dto"
)

type LogicInterface interface {
	GetBestAd() (*dto.AdDTO, error)
	StartTicker()
	BrakeAd(adId uint, duration time.Duration)
}

type LogicService struct {
	adsList       []*dto.AdDTO
	brakedAdIds   map[uint]struct{}
	getAdsAPIPath string
	interval      int
}

func NewLogicService() LogicInterface {
	interval, _ := strconv.Atoi(os.Getenv("ADS_FETCH_INTERVAL_SECS"))
	return &LogicService{
		adsList:       make([]*dto.AdDTO, 0),
		brakedAdIds:   make(map[uint]struct{}, 0),
		getAdsAPIPath: "http://" + os.Getenv("PANEL_HOSTNAME") + ":" + os.Getenv("PANEL_PORT") + os.Getenv("GET_ADS_API"),
		interval:      interval,
	}
}

func (ls *LogicService) isBetterThan(lhs, rhs *dto.AdDTO) bool {
	return rhs.Score > lhs.Score
}

func (ls *LogicService) GetBestAd() (*dto.AdDTO, error) {
	if len(ls.adsList) == 0 {
		return nil, fmt.Errorf("no ad was found")
	}

	index := 0
	for index = 0; index < len(ls.adsList); index++ {
		_, ok := ls.brakedAdIds[ls.adsList[index].ID]
		if !ok {
			break
		}
	}

	if index == len(ls.adsList) {
		return nil, fmt.Errorf("no ad was found")
	}
	bestAd := ls.adsList[index]
	anyValidMap := false

	for i := index; i < len(ls.adsList); i++ {
		ad := ls.adsList[i]
		_, ok := ls.brakedAdIds[ad.ID]
		if ok {
			continue
		}
		anyValidMap = true

		if ls.isBetterThan(bestAd, ad) {
			bestAd = ad
		}
	}

	if anyValidMap {
		return bestAd, nil
	} else {
		return nil, fmt.Errorf("no ad was found")
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

	var newAdsList []*dto.AdDTO
	for _, ad := range ads {
		newAdsList = append(newAdsList, &ad)
	}

	fmt.Printf("%d new ads were fetched.\n", len(newAdsList)-len(ls.adsList))
	ls.adsList = newAdsList

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

func (ls *LogicService) BrakeAd(adId uint, duration time.Duration) {
	ls.brakedAdIds[adId] = struct{}{}

	// Start a new goroutine to remove the adId after the specified duration
	go func() {
		time.Sleep(duration)
		delete(ls.brakedAdIds, adId)
	}()
}
