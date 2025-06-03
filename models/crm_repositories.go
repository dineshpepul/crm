package models

// CRMRepositories contains all CRM related repositories
type CRMRepositories struct {
	LeadRepo            LeadRepository
	LeadFieldConfigRepo LeadFieldConfigRepository
	ContactRepo         ContactRepository
	DealRepo            DealRepository
	CampaignRepo        CampaignRepository
	DashboardRepo       DashboardRepository
	AnalyticsRepo       AnalyticsRepository
	TargetRepo          TargetRepository
	NurtureRepo         NurtureRepository
	UserRepo            UserRepository
	LeadScoreType       ScoreRepository
}
