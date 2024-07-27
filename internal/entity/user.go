package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"uniqueIndex"`
	Password string `json:"password,omitempty"`
}
