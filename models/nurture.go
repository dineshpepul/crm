
package models

import (
	"time"
	
	"gorm.io/gorm"
)

// NurtureSequence represents an automated sequence of nurturing activities
type NurtureSequence struct {
	ID          int            `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Description string         `json:"description" gorm:"type:text"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Steps       []NurtureStep  `json:"steps" gorm:"foreignKey:SequenceID"`
}

// NurtureStep represents a single step in a nurture sequence
type NurtureStep struct {
	ID          int            `json:"id" gorm:"primaryKey"`
	SequenceID  int            `json:"sequence_id" gorm:"not null"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Type        string         `json:"type" gorm:"size:50;not null"` // email, task, notification, etc.
	Content     string         `json:"content" gorm:"type:text"`
	Delay       int            `json:"delay" gorm:"default:0"`       // delay in hours from previous step
	OrderIndex  int            `json:"order_index" gorm:"not null"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	Conditions  string         `json:"conditions" gorm:"type:text"`  // JSON for conditions
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// NurtureEnrollment represents a lead enrolled in a nurture sequence
type NurtureEnrollment struct {
	ID         int            `json:"id" gorm:"primaryKey"`
	SequenceID int            `json:"sequence_id" gorm:"not null"`
	LeadID     int            `json:"lead_id" gorm:"not null"`
	Status     string         `json:"status" gorm:"size:50;default:'active'"`
	CurrentStep int           `json:"current_step"`
	StartedAt  time.Time      `json:"started_at"`
	CompletedAt *time.Time    `json:"completed_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// NurtureActivity represents activity within a nurture sequence
type NurtureActivity struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	EnrollmentID int       `json:"enrollment_id" gorm:"not null"`
	StepID       int       `json:"step_id" gorm:"not null"`
	Type         string    `json:"type" gorm:"size:50;not null"` // sent, opened, clicked, completed, etc.
	Details      string    `json:"details" gorm:"type:text"`     // JSON for additional details
	CreatedAt    time.Time `json:"created_at"`
}
