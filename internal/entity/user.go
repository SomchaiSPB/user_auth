package entity

import (
	"time"
)

// User represents a user entity
type User struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Name      string     `json:"name" gorm:"uniqueIndex"`
	Password  string     `json:"password,omitempty"`
}
