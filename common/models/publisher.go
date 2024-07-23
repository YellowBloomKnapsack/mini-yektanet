package models

import "gorm.io/gorm"

type Publisher struct {
	gorm.Model
	Username       string `gorm:"not null;index"`
	Balance        int64  `gorm:"not null;default:0"`
	AdsInteraction []AdsInteraction
	Transactions    []Transaction `gorm:"polymorphicType:CustomerType;polymorphicId:CustomerID;polymorphicValue:1"`
}
