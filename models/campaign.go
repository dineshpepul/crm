
package models

import (
	"time"
)

// Campaign represents a marketing campaign in the CRM system
type Campaign struct {
	ID           int        `json:"id" gorm:"primaryKey"`
	Name         string     `json:"name" gorm:"size:255;not null"`
	Description  string     `json:"description" gorm:"type:text"`
	CampaignType string     `json:"campaign_type" gorm:"size:50;not null"`
	Status       string     `json:"status" gorm:"size:50;not null;default:'draft'"`
	StartDate    *time.Time `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	Budget       float64    `json:"budget"`
	Currency     string     `json:"currency" gorm:"size:3;default:'USD'"`
	CreatedBy    int        `json:"created_by" gorm:"not null"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// CampaignLead represents the relationship between campaigns and leads
type CampaignLead struct {
	CampaignID int       `json:"campaign_id" gorm:"primaryKey"`
	LeadID     int       `json:"lead_id" gorm:"primaryKey"`
	Status     string    `json:"status" gorm:"size:50;default:'active'"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CampaignTemplate represents an email template for campaigns
type CampaignTemplate struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"size:255;not null"`
	Subject      string    `json:"subject" gorm:"size:255;not null"`
	Content      string    `json:"content" gorm:"type:text;not null"`
	TemplateType string    `json:"template_type" gorm:"size:50;not null"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	CreatedBy    int       `json:"created_by" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
