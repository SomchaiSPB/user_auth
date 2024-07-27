package entity

import (
	"time"
)

// Product represents a product entity
type Product struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	Name        string     `json:"name" gorm:"uniqueIndex"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
}
