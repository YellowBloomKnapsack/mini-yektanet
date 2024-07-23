package models

import (
	"time"

	"gorm.io/gorm"
)

type AdsInteraction struct {
	gorm.Model
	// Count    int   `gorm:"not null;default:0"` //count at that period of time
	Type        int `gorm:"not null"`
	EventTime   time.Time
	Bid         int64 `gorm:"not null"`
	AdID        uint
	Ad          Ad
	PublisherID uint
	Publisher   Publisher
}

type AdsInteractionType int

const (
	Impression AdsInteractionType = 0
	Click      AdsInteractionType = 1
)
