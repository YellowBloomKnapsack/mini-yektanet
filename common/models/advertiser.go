package models

import "gorm.io/gorm"

type Advertiser struct {
	gorm.Model
	Username string `gorm:"not null;index"`
	Balance  int64  `gorm:"not null;default:0"`
	Ads 	[]Ad
}
