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

var adsList = make([]*dto.AdDTO, 0)

func isBetterThan(lhs, rhs *dto.AdDTO) bool {
	return rhs.Bid > lhs.Bid
}

func GetBestAd() (*dto.AdDTO, error) {
	if len(adsList) == 0 {
		return nil, fmt.Errorf("no ad was found")
	}

	bestAd := adsList[0]

	for _, ad := range adsList {
		if isBetterThan(bestAd, ad) {
			bestAd = ad
		}
	}

	return bestAd, nil
}

func updateAdsList() error {
	fmt.Println("Fetching ads from panel API...")

	getAdsAPIPath := "http://" + os.Getenv("HOSTNAME") + ":" + os.Getenv("PANEL_PORT") + os.Getenv("GET_ADS_API")
	resp, err := http.Get(getAdsAPIPath)
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

	fmt.Printf("%d new ads were fetched.\n", len(newAdsList) - len(adsList))
	adsList = newAdsList

	return nil
}

func StartTicker() {
	interval, _ := strconv.Atoi(os.Getenv("ADS_FETCH_INTERVAL_SECS"))
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := updateAdsList()
			if err != nil {
			    fmt.Println(err)
			}
		}
	}
}
