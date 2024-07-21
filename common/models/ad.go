package models

import "gorm.io/gorm"

type Ad struct {
	gorm.Model
	Text            string
	ImagePath       string
	Bid             int64 `gorm:"not null"`
	Active          bool  `gorm:"default:true;index"`
	TotalCost       int64 `gorm:"default:0"`
	Website         string
	AdvertiserID    uint
	Advertiser      Advertiser
	AdsInteractions []AdsInteraction
}
