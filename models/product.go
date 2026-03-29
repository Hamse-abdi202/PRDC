package models

import "gorm.io/gorm"

type Product struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name" gorm:"not null"`
	SKU           string  `json:"sku" gorm:"unique;not null"`
	PurchasePrice float64 `json:"purchaseprice" gorm:"not null"`
	SellingPrice  float64 `json:"sellingprice" gorm:"not null"`
	Quantity      int     `json:"quantity" gorm:"not null"`
	CategoryID    uint
	gorm.Model
}
