package models

import (
	"time"
)

// LeadRepository interface for lead operations
type LeadRepository interface {
	FindByID(id int) (*Lead, error)
	List(companyId int) ([]Lead, error)
	ListByStatus(status string) ([]Lead, error)
	ListByAssignee(assigneeID int) ([]Lead, error)
	Create(lead []CrmFieldData) error
	Update(lead *Lead) error
	Delete(id int) error
	ValidateLeadFields(lead *Lead, requiredFields []string) error
	GetLastSubmitId() (int, error)
}

// LeadFieldConfigRepository interface for lead field configuration
type LeadFieldConfigRepository interface {
	GetAllFieldConfigs(companyId int) ([]LeadFieldConfig, error)
	GetVisibleFieldConfigs(companyId int) ([]LeadFieldConfig, error)
	GetRequiredFieldConfigs(companyId int) ([]LeadFieldConfig, error)
	GetFieldConfigsBySection(section string) ([]LeadFieldConfig, error)
	GetFieldConfig(id int) (*LeadFieldConfig, error)
	CreateFieldConfig(config *LeadFieldConfig) error
	UpdateFieldConfig(config *LeadFieldConfig) error
	DeleteFieldConfig(id uint) error
	ReorderFormFields(fieldIDs []uint) error
	GetAllFormSections(companyId int) ([]LeadFormSection, error)
	GetVisibleFormSections(companyId int) ([]LeadFormSection, error)
	CreateFormSection(section *LeadFormSection) error
	UpdateFormSection(section *LeadFormSection) error
	DeleteFormSection(id int) error
	ReorderFormSections(sectionIDs []int) error
	GetFormStructure(companyId int) (map[string]interface{}, error)
}

type ScoreRepository interface {
	ScoreUpdateRepo(config []ScoreType) error
}

// ContactRepository interface for contact operations
type ContactRepository interface {
	FindByID(id int) (*Contact, error)
	List(offset int, limit int, companyId int) ([]Contact, error)
	FindByLead(leadID int) ([]Contact, error)
	Create(contact *Contact) error
	Update(contact *Contact) error
	Delete(id int) error
	Search(query string, companyId int) ([]Contact, error)
}

// DealRepository interface for deal operations
type DealRepository interface {
	FindByID(id int) (*Deal, error)
	List(offset int, limit int, filters map[string]interface{}, companyId int) ([]Deal, error)
	FindByLead(leadID int) ([]Deal, error)
	Create(deal *Deal) error
	Update(deal *Deal) error
	Delete(id int) error
	GetDealPipeline(companyId int) ([]map[string]interface{}, error)
}

// CampaignRepository interface for campaign operations
type CampaignRepository interface {
	GetCampaigns(offset int, limit int) ([]Campaign, error)
	GetCampaignByID(id int) (*Campaign, error)
	CreateCampaign(campaign *Campaign) error
	UpdateCampaign(campaign *Campaign) error
	DeleteCampaign(id int) error
	GetCampaignStats(id int) (map[string]interface{}, error)
	GetLeadsForCampaign(id int) ([]Lead, error)
	AssignLeadsToCampaign(campaignID int, leadIDs []int) error
	RemoveLeadsFromCampaign(campaignID int, leadIDs []int) error
	GetTemplates(offset int, limit int) ([]CampaignTemplate, error)
	GetTemplateByID(id int) (*CampaignTemplate, error)
	CreateTemplate(template *CampaignTemplate) error
	UpdateTemplate(template *CampaignTemplate) error
	DeleteTemplate(id int) error
}

// DashboardRepository interface for dashboard operations
type DashboardRepository interface {
	GetDashboardSummary(companyId int) (map[string]interface{}, error)
	GetLeadsBySource(companyId int) ([]map[string]interface{}, error)
	GetLeadsByStatus(companyId int) ([]map[string]interface{}, error)
	GetRevenueByMonth(year int, companyId int) ([]map[string]interface{}, error)
	GetSalesForecast(months int, companyId int) ([]map[string]interface{}, error)
}

// AnalyticsRepository interface for analytics operations
type AnalyticsRepository interface {
	GetLeadAnalytics(startDate time.Time, endDate time.Time, companyId int) (map[string]interface{}, error)
	GetDealAnalytics(startDate time.Time, endDate time.Time, companyId int) (map[string]interface{}, error)
	GetSalesActivity(startDate time.Time, endDate time.Time, companyId int) (map[string]interface{}, error)
	GetPerformanceByUser(startDate time.Time, endDate time.Time, companyId int) (map[string]interface{}, error)
	GetFunnelAnalytics(companyId int) (map[string]interface{}, error)
	GetTargetAnalytics(startDate time.Time, endDate time.Time, userId *uint, companyId int) (map[string]interface{}, error)
}

// TargetRepository interface for sales target operations
type TargetRepository interface {
	GetTargets(filters map[string]interface{}, companyId int) ([]Target, error)
	GetTargetByID(id int) (*Target, error)
	CreateTarget(target *Target) error
	UpdateTarget(target *Target) error
	DeleteTarget(id int) error
	GetTargetProgress(id int, companyId int) (map[string]interface{}, error)
	GetAllTargetProgress(companyId int) ([]map[string]interface{}, error)
}

// NurtureRepository interface for nurture sequences
type NurtureRepository interface {
	GetSequences(offset int, limit int) ([]NurtureSequence, error)
	GetSequenceByID(id int) (*NurtureSequence, error)
	CreateSequence(sequence *NurtureSequence) error
	UpdateSequence(sequence *NurtureSequence) error
	DeleteSequence(id int) error
	GetStepsBySequence(sequenceID int) ([]NurtureStep, error)
	CreateStep(step *NurtureStep) error
	UpdateStep(step *NurtureStep) error
	DeleteStep(id int) error
	GetEnrollments(sequenceID int, offset int, limit int) ([]NurtureEnrollment, error)
	EnrollLead(enrollment *NurtureEnrollment) error
	UpdateEnrollment(enrollment *NurtureEnrollment) error
	GetEnrollmentActivity(enrollmentID int) ([]NurtureActivity, error)
	RecordActivity(activity *NurtureActivity) error

	// Campaign related methods
	GetCampaigns(offset int, limit int, companyId int) ([]Campaign, error)
	GetCampaignByID(id int) (*Campaign, error)
	CreateCampaign(campaign *Campaign) error
	UpdateCampaign(campaign *Campaign) error
	DeleteCampaign(id int) error
	GetCampaignStats(id int) (map[string]interface{}, error)
	GetLeadsForCampaign(id int) ([]Lead, error)
	AssignLeadsToCampaign(campaignID int, leadIDs []int) error
	RemoveLeadsFromCampaign(campaignID int, leadIDs []int) error
	GetTemplates(offset int, limit int, companyId int) ([]CampaignTemplate, error)
	GetTemplateByID(id int) (*CampaignTemplate, error)
	CreateTemplate(template *CampaignTemplate) error
	UpdateTemplate(template *CampaignTemplate) error
	DeleteTemplate(id int) error
}

// UserRepository interface for user operations
type UserRepository interface {
	FindByID(id int) (*User, error)
	FindByEmail(email string) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id int) error
	List() ([]User, error)
}
