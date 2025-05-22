
package models

import (
	"gorm.io/gorm"
)

// Repositories contains all CRM repositories
type Repositories struct {
	db *gorm.DB
	LeadRepo           LeadRepository
	LeadFieldConfigRepo LeadFieldConfigRepository
	DealRepo           DealRepository
	ContactRepo        ContactRepository
	DashboardRepo      DashboardRepository
	TargetRepo         TargetRepository
	NurtureRepo        NurtureRepository
	UserRepo           UserRepository
	CampaignRepo       CampaignRepository
	AnalyticsRepo      AnalyticsRepository
}

// NewRepositories initializes repositories
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		db: db,
	}
}
