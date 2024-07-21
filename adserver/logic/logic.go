package logic

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"YellowBloomKnapsack/mini-yektanet/common/dto"
)

var (
	tmpPrivateKey string
	privateKey    []byte
)

func Init() {
	tmpPrivateKey = os.Getenv("PRIVATE_KEY")
	privateKey, _ = base64.StdEncoding.DecodeString(tmpPrivateKey)
}

type InteractionType uint8

const (
	ImpressionType InteractionType = 0
	ClickType      InteractionType = 1
)

type CustomToken struct {
	Interaction       InteractionType `json:"interaction"`
	AdID              uint            `json:"ad_id"`
	PublisherUsername string          `json:"publisher_username"`
	RedirectPath      string          `json:"redirect_path"`
	CreatedAt         int64           `json:"created_at"`
}

func encrypt(data []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func generateToken(interaction InteractionType, adID uint, publisherUsername, redirectPath string, key []byte) (string, error) {
	token := CustomToken{
		Interaction:       interaction,
		AdID:              adID,
		PublisherUsername: publisherUsername,
		RedirectPath:      redirectPath,
		CreatedAt:         time.Now().Unix(),
	}

	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return "", err
	}

	encryptedToken, err := encrypt(tokenBytes, key)
	if err != nil {
		return "", err
	}

	return encryptedToken, nil
}

func GenerateToken(interaction InteractionType, adID uint, publisherUsername, redirectPath string) (string, error) {
	return generateToken(interaction, adID, publisherUsername, redirectPath, privateKey)
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

func updateAdsList() error {
	fmt.Println("Fetching ads from panel API")

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
			fmt.Println(err)
		}
	}
}
