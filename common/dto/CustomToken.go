package dto

import (
	"YellowBloomKnapsack/mini-yektanet/common/models"
)

type CustomToken struct {
	Interaction models.AdsInteractionType `json:"interaction"`
	AdID        uint                      `json:"ad_id"`
	PublisherID uint                      `json:"publisher_id"`
	Bid         int64                     `json:"bid"`
	CreatedAt   int64                     `json:"created_at"`
}
