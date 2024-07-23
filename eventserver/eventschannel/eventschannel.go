package eventschannel

import (
	"fmt"
	"os"
	"time"
	"strconv"
	"bytes"
	"encoding/json"
	"net/http"

	"YellowBloomKnapsack/mini-yektanet/common/dto"
)

type ClickEvent struct {
	Data dto.CustomToken
	ClickTime time.Time
}

type ImpressionEvent struct {
	Data dto.CustomToken
	ImpressionTime time.Time
}

var (
	clickEvents      chan ClickEvent
	impressionEvents chan ImpressionEvent 
	clickApiUrl      string
	impressionApiUrl string
)

func initialize() {
	baseUrl := "http://" + os.Getenv("HOSTNAME") + ":" + os.Getenv("PANEL_PORT") 
	clickBuffSize, _ := strconv.Atoi(os.Getenv("CLICK_BUFF_SIZE"))
	impressionBuffSize, _ := strconv.Atoi(os.Getenv("IMPRESSION_BUFF_SIZE"))

	clickEvents      = make(chan ClickEvent, clickBuffSize)
	impressionEvents = make(chan ImpressionEvent, impressionBuffSize)
	clickApiUrl      = baseUrl + os.Getenv("INTERACTION_CLICK_API")
	impressionApiUrl = baseUrl + os.Getenv("INTERACTION_IMPRESSION_API")
}

func Start() {
    initialize()
	go func() {
		for {
			select {
			case event := <-clickEvents:
				callClickApi(event)
			case event := <-impressionEvents:
				callImpressionApi(event)
			}
		}
	}()
}

func InvokeClickEvent(data *dto.CustomToken, clickTime time.Time) {
	event := ClickEvent{
		Data: *data,
		ClickTime: clickTime,
	}

	clickEvents <- event
}

func InvokeImpressionEvent(data *dto.CustomToken, impressionTime time.Time) {
	event := ImpressionEvent{
		Data: *data,
		ImpressionTime: impressionTime,
	}

	impressionEvents <- event
}

func callClickApi(event ClickEvent) {
	dataDto := dto.InteractionDto{
		PublisherUsername: event.Data.PublisherUsername,
		ClickTime:         event.ClickTime,
		AdID:              event.Data.AdID,
	}

	jsonData, err := json.Marshal(dataDto)
	if err != nil {
		fmt.Printf("failed to marshal InteractionDto to JSON: %v\n", err)
		return
	}

	resp, err := http.Post(clickApiUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("failed to send POST request: %v\n", err)
		// reinsert event to clickEvents if panel is down
		clickEvents <- event
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("failed to fetch ads: status code %d\n", resp.StatusCode)
	}
}

func callImpressionApi(event ImpressionEvent) {
	dataDto := dto.InteractionDto{
		PublisherUsername: event.Data.PublisherUsername,
		ClickTime:         event.ImpressionTime,
		AdID:              event.Data.AdID,
	}

	jsonData, err := json.Marshal(dataDto)
	if err != nil {
		fmt.Printf("failed to marshal InteractionDto to JSON: %v\n", err)
		return
	}

	resp, err := http.Post(impressionApiUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("failed to send POST request: %v\n", err)
		// reinsert event to impressionEvents if panel is down
		impressionEvents <- event
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("failed to fetch ads: status code %d\n", resp.StatusCode)
	}
}
