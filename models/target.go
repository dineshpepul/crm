
package models

import (
	"time"

	"gorm.io/gorm"
)

// Target represents a sales target in the CRM system
type Target struct {
	ID             int            `json:"id" gorm:"primaryKey"`
	Name           string         `json:"name" gorm:"size:100;not null"`
	TargetType     string         `json:"target_type" gorm:"size:50;not null"` // revenue, leads, deals, etc.
	TargetValue    float64        `json:"target_value" gorm:"not null"`
	ActualValue    float64        `json:"actual_value"`
	UserId         *int           `json:"user_id"` // If assigned to a specific user
	TeamId         *int           `json:"team_id"` // If assigned to a team
	StartDate      time.Time      `json:"start_date" gorm:"not null"`
	EndDate        time.Time      `json:"end_date" gorm:"not null"`
	Period         string         `json:"period" gorm:"size:50;not null"` // monthly, quarterly, annual
	Status         string         `json:"status" gorm:"size:50;default:'active'"`
	Currency       string         `json:"currency" gorm:"size:3;default:'USD'"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
