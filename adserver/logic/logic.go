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
	"YellowBloomKnapsack/mini-yektanet/common/models"
)

const (
	getAdsAPIPath = "localhsot:8082/ads" // TODO: make this an env var in panel
)

func ToAdDTO(dbAd models.Ad) *dto.AdDTO {
	return &dto.AdDTO{
		ID:        dbAd.ID,
		Text:      dbAd.Text,
		ImagePath: dbAd.ImagePath,
		Bid:       dbAd.Bid,
		Website:   dbAd.Website,
	}
}

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

// WARNING: not tested yet
func updateAdsList() error {
	adsList = append(adsList, &dto.AdDTO{ID: 1, Text: "salam", ImagePath: "somewhere", Bid: 30, Website: "google.com"})
	adsList = append(adsList, &dto.AdDTO{ID: 2, Text: "khodafez", ImagePath: "somewhere but not here", Bid: 40, Website: "duckduckgo.com"})
	fmt.Println("Yayyyyyyyyyy!")
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

	var ads []models.Ad
	if err := json.Unmarshal(body, &ads); err != nil {
		return fmt.Errorf("failed to unmarshal ads: %v", err)
	}

	var newAdsList []*dto.AdDTO
	for _, ad := range ads {
		newAdsList = append(newAdsList, ToAdDTO(ad))
	}

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
			_ = updateAdsList()
		}
	}
}
