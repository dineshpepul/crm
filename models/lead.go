package models

import (
	"time"

	"gorm.io/gorm"
)

// Lead represents a lead in the CRM system
type Lead struct {
	ID           uint              `json:"id" gorm:"primaryKey"`
	Name         string            `json:"name" gorm:"size:255;not null"`
	Email        string            `json:"email" gorm:"size:255"`
	Phone        string            `json:"phone" gorm:"size:50"`
	Company      string            `json:"company" gorm:"size:255"`
	Source       string            `json:"source" gorm:"size:100"`
	Status       string            `json:"status" gorm:"size:50;not null;default:'new'"`
	Score        *int              `json:"score" gorm:"default:null"`
	AssignedToID *uint             `json:"assigned_to_id" gorm:"default:null"`
	AssignedTo   *User             `json:"assigned_to" gorm:"foreignKey:AssignedToID"` // Added field for relationships
	Notes        string            `json:"notes" gorm:"type:text"`
	Tags         []string          `json:"tags" gorm:"-"` // Handled through a separate table
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	DeletedAt    gorm.DeletedAt    `json:"deleted_at" gorm:"index"`
	CustomFields []LeadCustomField `json:"custom_fields" gorm:"foreignKey:LeadID"`
	Type         string            `json:"type" gorm:"default:null"`
}

// LeadTag represents a tag associated with a lead
type LeadTag struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	LeadID    uint      `json:"lead_id"`
	Tag       string    `json:"tag" gorm:"size:50"`
	CreatedAt time.Time `json:"created_at"`
}

type ResultLead struct {
	ID           uint              `json:"id" gorm:"primaryKey"`
	Name         string            `json:"name" gorm:"size:255;not null"`
	Email        string            `json:"email" gorm:"size:255"`
	Phone        string            `json:"phone" gorm:"size:50"`
	Company      string            `json:"company" gorm:"size:255"`
	Source       string            `json:"source" gorm:"size:100"`
	Status       string            `json:"status" gorm:"size:50;not null;default:'new'"`
	Score        *int              `json:"score" gorm:"default:null"`
	AssignedToID *uint             `json:"assigned_to_id" gorm:"default:null"`
	AssignedTo   *User             `json:"assigned_to" gorm:"foreignKey:AssignedToID"` // Added field for relationships
	Notes        string            `json:"notes" gorm:"type:text"`
	Tags         []string          `json:"tags" gorm:"-"` // Handled through a separate table
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	DeletedAt    gorm.DeletedAt    `json:"deleted_at" gorm:"index"`
	CustomFields []LeadCustomField `json:"custom_fields" gorm:"foreignKey:LeadID"`
	Type         string            `json:"type"`
}
