package repositories

import (
	"crm-app/backend/models"

	"gorm.io/gorm"
)

// NewRepositoriesInit initializes all repositories
func NewRepositoriesInit(db *gorm.DB) *models.Repositories {
	repos := models.NewRepositories(db)
	repos.LeadRepo = NewLeadRepository(db)
	repos.LeadFieldConfigRepo = NewLeadFieldConfigRepository(db)
	repos.DealRepo = NewDealRepository(db)
	repos.ContactRepo = NewContactRepository(db)
	repos.DashboardRepo = NewDashboardRepository(db)
	repos.TargetRepo = NewTargetRepository(db)
	repos.NurtureRepo = NewNurtureRepository(db)
	repos.UserRepo = NewUserRepository(db)
	repos.CampaignRepo = NewCampaignRepository(db)
	repos.AnalyticsRepo = NewAnalyticsRepository(db)

	return repos
}

// NewCRMRepositories initializes all repositories
func NewCRMRepositories(db *gorm.DB) *models.CRMRepositories {
	return &models.CRMRepositories{
		LeadRepo:            NewLeadRepository(db),
		LeadFieldConfigRepo: NewLeadFieldConfigRepository(db),
		ContactRepo:         NewContactRepository(db),
		DealRepo:            NewDealRepository(db),
		CampaignRepo:        NewCampaignRepository(db),
		DashboardRepo:       NewDashboardRepository(db),
		AnalyticsRepo:       NewAnalyticsRepository(db),
		TargetRepo:          NewTargetRepository(db),
		NurtureRepo:         NewNurtureRepository(db),
		UserRepo:            NewUserRepository(db),
	}
}

// GormContactRepository is already defined in contact_repository.go

type gormLeadRepository struct {
	db *gorm.DB
}

type gormDealRepository struct {
	db *gorm.DB
}

type gormCampaignRepository struct {
	db *gorm.DB
}

type gormDashboardRepository struct {
	db *gorm.DB
}

type gormAnalyticsRepository struct {
	db *gorm.DB
}

type gormTargetRepository struct {
	db *gorm.DB
}

type gormNurtureRepository struct {
	db *gorm.DB
}

type gormUserRepository struct {
	db *gorm.DB
}

// NewLeadRepository creates a new lead repository
func NewLeadRepository(db *gorm.DB) models.LeadRepository {
	return &gormLeadRepository{db: db}
}

// NewLeadFieldConfigRepository creates a new lead field config repository
func NewLeadFieldConfigRepository(db *gorm.DB) models.LeadFieldConfigRepository {
	return &GormLeadFieldConfigRepository{db: db}
}

// NewContactRepository creates a new contact repository
func NewContactRepository(db *gorm.DB) models.ContactRepository {
	return &GormContactRepository{db: db}
}

// NewDealRepository creates a new deal repository
func NewDealRepository(db *gorm.DB) models.DealRepository {
	return &gormDealRepository{db: db}
}

// NewCampaignRepository creates a new campaign repository
func NewCampaignRepository(db *gorm.DB) models.CampaignRepository {
	return &gormCampaignRepository{db: db}
}

// NewDashboardRepository creates a new dashboard repository
func NewDashboardRepository(db *gorm.DB) models.DashboardRepository {
	return &gormDashboardRepository{db: db}
}

// NewAnalyticsRepository creates a new analytics repository
func NewAnalyticsRepository(db *gorm.DB) models.AnalyticsRepository {
	return &gormAnalyticsRepository{db: db}
}

// NewTargetRepository creates a new target repository
func NewTargetRepository(db *gorm.DB) models.TargetRepository {
	return &gormTargetRepository{db: db}
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) models.UserRepository {
	return &gormUserRepository{db: db}
}

// NewNurtureRepository creates a new nurture repository
func NewNurtureRepository(db *gorm.DB) models.NurtureRepository {
	return &gormNurtureRepository{db: db}
}
