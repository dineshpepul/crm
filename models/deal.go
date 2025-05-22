
package models

import (
	"time"

	"gorm.io/gorm"
)

// Deal represents a deal in the CRM system
type Deal struct {
	ID                int            `json:"id" gorm:"primaryKey"`
	LeadID            int            `json:"lead_id" gorm:"not null"`
	Title             string         `json:"title" gorm:"size:255;not null"`
	Amount            float64        `json:"amount"`
	Currency          string         `json:"currency" gorm:"size:3;default:'USD'"`
	Stage             string         `json:"stage" gorm:"size:50;not null"`
	Probability       int            `json:"probability"` // 0-100 percent
	ExpectedCloseDate *time.Time     `json:"expected_close_date"`
	AssignedTo        *int           `json:"assigned_to"`
	Notes             string         `json:"notes,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Tags              []string       `json:"tags,omitempty" gorm:"-"`
}
