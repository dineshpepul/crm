
package models

import (
	"time"

	"gorm.io/gorm"
)

// Contact represents a contact in the CRM system
type Contact struct {
	ID         int            `json:"id" gorm:"primaryKey"`
	LeadID     *int           `json:"lead_id" gorm:"index"`
	Name       string         `json:"name" gorm:"size:255;not null"`
	Email      string         `json:"email" gorm:"size:255"`
	Phone      string         `json:"phone" gorm:"size:50"`
	Position   string         `json:"position" gorm:"size:100"`
	IsPrimary  bool           `json:"is_primary" gorm:"default:false"`
	Notes      string         `json:"notes,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
