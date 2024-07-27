package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"YellowBloomKnapsack/mini-yektanet/common/dto"
)

type WorkerInterface interface {
	Start()
	InvokeClickEvent(data *dto.CustomToken, clickTime time.Time)
	InvokeImpressionEvent(data *dto.CustomToken, impressionTime time.Time)
}

type WorkerService struct {
	clickEvents        chan ClickEvent
	impressionEvents   chan ImpressionEvent
	clickApiUrl        string
	impressionApiUrl   string
	clickBuffSize      int
	impressionBuffSize int
}

func NewWorkerService() WorkerInterface {
	baseUrl := "http://" + os.Getenv("PANEL_HOSTNAME") + ":" + os.Getenv("PANEL_PORT")
	clickBuffSize, _ := strconv.Atoi(os.Getenv("CLICK_BUFF_SIZE"))
	impressionBuffSize, _ := strconv.Atoi(os.Getenv("IMPRESSION_BUFF_SIZE"))

	return &WorkerService{
		clickEvents:        make(chan ClickEvent, clickBuffSize),
		impressionEvents:   make(chan ImpressionEvent, impressionBuffSize),
		clickApiUrl:        baseUrl + os.Getenv("INTERACTION_CLICK_API"),
		impressionApiUrl:   baseUrl + os.Getenv("INTERACTION_IMPRESSION_API"),
		clickBuffSize:      clickBuffSize,
		impressionBuffSize: impressionBuffSize,
	}
}

type ClickEvent struct {
	Data      dto.CustomToken
	ClickTime time.Time
}

type ImpressionEvent struct {
	Data           dto.CustomToken
	ImpressionTime time.Time
}

func (w *WorkerService) Start() {
	go func() {
		for {
			select {
			case event := <-w.clickEvents:
				w.callClickApi(event)
			case event := <-w.impressionEvents:
				w.callImpressionApi(event)
			}
		}
	}()
}

func (w *WorkerService) InvokeClickEvent(data *dto.CustomToken, clickTime time.Time) {
	event := ClickEvent{
		Data:      *data,
		ClickTime: clickTime,
	}

	w.clickEvents <- event
}

func (w *WorkerService) InvokeImpressionEvent(data *dto.CustomToken, impressionTime time.Time) {
	event := ImpressionEvent{
		Data:           *data,
		ImpressionTime: impressionTime,
	}

	w.impressionEvents <- event
}

func (w *WorkerService) callClickApi(event ClickEvent) {
	dataDto := dto.InteractionDto{
		PublisherUsername: event.Data.PublisherUsername,
		EventTime:         event.ClickTime,
		AdID:              event.Data.AdID,
	}

	jsonData, err := json.Marshal(dataDto)
	if err != nil {
		fmt.Printf("failed to marshal InteractionDto to JSON: %v\n", err)
		return
	}

	resp, err := http.Post(w.clickApiUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("failed to send POST request: %v\n", err)
		// reinsert event to clickEvents if panel is down
		w.clickEvents <- event
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("failed to fetch ads: status code %d\n", resp.StatusCode)
	}
}

func (w *WorkerService) callImpressionApi(event ImpressionEvent) {
	dataDto := dto.InteractionDto{
		PublisherUsername: event.Data.PublisherUsername,
		EventTime:         event.ImpressionTime,
		AdID:              event.Data.AdID,
	}

	jsonData, err := json.Marshal(dataDto)
	if err != nil {
		fmt.Printf("failed to marshal InteractionDto to JSON: %v\n", err)
		return
	}

	resp, err := http.Post(w.impressionApiUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("failed to send POST request: %v\n", err)
		// reinsert event to impressionEvents if panel is down
		w.impressionEvents <- event
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("failed to fetch ads: status code %d\n", resp.StatusCode)
	}
}
