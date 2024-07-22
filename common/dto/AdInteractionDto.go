package dto

import "time"
type InteractionDto struct {
    PublisherUsername string    `json:"publisherUsername" binding:"required"`
    ClickTime         time.Time `json:"clickTime" binding:"required"`
    AdID              uint      `json:"adId" binding:"required"`
}