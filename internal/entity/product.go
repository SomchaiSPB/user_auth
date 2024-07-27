package entity

import "gorm.io/gorm"

// Product represents a product entity
type Product struct {
	gorm.Model
	Name        string  `json:"name" gorm:"uniqueIndex"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
