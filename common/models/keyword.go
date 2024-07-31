package models

import "gorm.io/gorm"

type Keyword struct {
	gorm.Model
	AdID     uint   `gorm:"index"`
	Keywords string `gorm:"type:varchar(255)"`
}
