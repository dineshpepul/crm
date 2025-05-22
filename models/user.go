package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID int `json:"id" gorm:"primaryKey"`
	// Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	Email        string    `json:"email" gorm:"type:varchar(255);unique;not null"`
	Name         string    `json:"name" gorm:"not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	Role         string    `json:"role" gorm:"default:user"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
