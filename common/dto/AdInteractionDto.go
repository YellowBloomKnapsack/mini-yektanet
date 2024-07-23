package dto

import "time"

type InteractionDto struct {
	PublisherUsername string    `json:"publisherUsername" binding:"required"`
	EventTime         time.Time `json:"clickTime" binding:"required"`
	AdID              uint      `json:"adId" binding:"required"`
}
