package entity

import "gorm.io/gorm"

// User represents a user entity
type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"uniqueIndex"`
	Password string `json:"password,omitempty"`
}
