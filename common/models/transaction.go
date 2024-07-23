package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	Description  string
	Amount       int64     `gorm:"not null"`
	Successful   bool      `gorm:"not null"`
	Income       bool      `gorm:"not null;index"`
	Time         time.Time `gorm:"not null"`
	CustomerID   uint
	CustomerType Customer
}

type Customer int

const (
	Customer_Advertiser Customer = 0
	Customer_Publisher  Customer = 1
)
