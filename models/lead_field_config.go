package models

import (
	"time"
)

// LeadFieldConfig defines customizable form fields for leads
type LeadFieldConfig struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	FieldName     string    `json:"field_name" gorm:"size:50;uniqueIndex;not null"`
	DisplayName   string    `json:"display_name" gorm:"size:100;not null"`
	FieldType     string    `json:"field_type" gorm:"size:20;not null"` // text, textarea, select, checkbox, etc.
	DefaultValue  string    `json:"default_value" gorm:"size:255"`
	Options       string    `json:"options" gorm:"type:text"` // JSON string for select options
	Required      bool      `json:"required" gorm:"default:false"`
	Visible       bool      `json:"visible" gorm:"default:true"`
	Section       string    `json:"section" gorm:"size:50;default:'default'"` // Which form section this belongs to
	OrderIndex    int       `json:"order_index" gorm:"not null;default:0"`
	HelpText      string    `json:"help_text" gorm:"size:255"`
	Placeholder   string    `json:"placeholder" gorm:"size:100"`
	ValidationMsg string    `json:"validation_msg" gorm:"size:255"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// LeadFormSection defines sections in the lead form
type LeadFormSection struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:50;uniqueIndex;not null"`
	Label       string    `json:"label" gorm:"size:100;not null"`
	Description string    `json:"description" gorm:"size:255"`
	OrderIndex  int       `json:"order_index" gorm:"not null;default:0"`
	Visible     bool      `json:"visible" gorm:"default:true"`
	Collapsible bool      `json:"collapsible" gorm:"default:false"`
	Expanded    bool      `json:"expanded" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// LeadCustomField stores custom field values for leads
type LeadCustomField struct {
	// ID         uint      `json:"id" gorm:"primaryKey"`
	LeadID     uint      `json:"lead_id" gorm:"not null"`
	FieldID    uint      `json:"field_id" gorm:"not null"`
	FieldName  string    `json:"field_name" gorm:"size:50;not null"` // For convenience and data integrity
	Value      string    `json:"value" gorm:"type:text"`
	FieldValue string    `json:"field_value" gorm:"type:text"` // Added field to match usage in code
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// type LeadScoreType struct {
// 	ID       int    `json:"id" gorm:"primaryKey"`
// 	Type     string `json:"type" gorm:"null"`
// 	MinScore int    `json:"type" gorm:"null"`
// 	MaxScore int    `json:"type" gorm:"null"`
// }
