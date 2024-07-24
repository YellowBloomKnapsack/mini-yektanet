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
}

type LogicService struct {
	adsList []*dto.AdDTO
	getAdsAPIPath string
	interval int
}

func NewLogicService() LogicInterface {
	interval, _ := strconv.Atoi(os.Getenv("ADS_FETCH_INTERVAL_SECS"))
	return &LogicService{
		adsList: make([]*dto.AdDTO, 0),
		getAdsAPIPath: "http://" + os.Getenv("HOSTNAME") + ":" + os.Getenv("PANEL_PORT") + os.Getenv("GET_ADS_API"),
		interval: interval,
	}
}

func (ls *LogicService) isBetterThan(lhs, rhs *dto.AdDTO) bool {
	return rhs.Bid > lhs.Bid
}

func (ls *LogicService) GetBestAd() (*dto.AdDTO, error) {
	if len(ls.adsList) == 0 {
		return nil, fmt.Errorf("no ad was found")
	}

	bestAd := ls.adsList[0]

	for _, ad := range ls.adsList {
		if ls.isBetterThan(bestAd, ad) {
			bestAd = ad
		}
	}

	return bestAd, nil
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

	fmt.Printf("%d new ads were fetched.\n", len(newAdsList) - len(ls.adsList))
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
