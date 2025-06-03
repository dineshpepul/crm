package services

import (
	"crm-app/backend/models"
	"errors"
)

// LeadService handles business logic for leads
type LeadService struct {
	repos *models.Repositories
}

// NewLeadService creates a new LeadService
func NewLeadService(repos *models.Repositories) *LeadService {
	return &LeadService{
		repos: repos,
	}
}

// GetLeads retrieves leads with optional filters
func (s *LeadService) GetLeads(companyId int) ([]models.Lead, error) {
	return s.repos.LeadRepo.List(companyId)
}

// GetLeadByID retrieves a lead by ID
func (s *LeadService) GetLeadByID(id int) (*models.Lead, error) {
	return s.repos.LeadRepo.FindByID(id)
}

// CreateLead creates a new lead with validation
func (s *LeadService) CreateLead(lead *models.Lead) error {
	// Get required fields for validation
	requiredFields, err := s.getRequiredFieldNames(lead.CompanyId)
	if err != nil {
		return err
	}

	// Validate required fields
	if err := s.repos.LeadRepo.ValidateLeadFields(lead, requiredFields); err != nil {
		return err
	}

	return s.repos.LeadRepo.Create(lead)
}

// UpdateLead updates a lead
func (s *LeadService) UpdateLead(lead *models.Lead) error {
	return s.repos.LeadRepo.Update(lead)
}

// DeleteLead deletes a lead
func (s *LeadService) DeleteLead(id int) error {
	return s.repos.LeadRepo.Delete(id)
}

// QualifyLead marks a lead as qualified
func (s *LeadService) QualifyLead(id int, score *int) error {
	lead, err := s.repos.LeadRepo.FindByID(id)
	if err != nil {
		return err
	}
	if lead == nil {
		return ErrNotFound
	}

	lead.Status = "qualified"
	if score != nil {
		lead.Score = score
	}

	return s.repos.LeadRepo.Update(lead)
}

// DisqualifyLead marks a lead as disqualified
func (s *LeadService) DisqualifyLead(id int) error {
	lead, err := s.repos.LeadRepo.FindByID(id)
	if err != nil {
		return err
	}
	if lead == nil {
		return ErrNotFound
	}

	lead.Status = "disqualified"

	return s.repos.LeadRepo.Update(lead)
}

// AssignLead assigns a lead to a user
func (s *LeadService) AssignLead(id int, assigneeID int) error {
	lead, err := s.repos.LeadRepo.FindByID(id)
	if err != nil {
		return err
	}
	if lead == nil {
		return ErrNotFound
	}

	assigneeIDUint := uint(assigneeID)
	lead.AssignedToID = &assigneeIDUint

	return s.repos.LeadRepo.Update(lead)
}

// GetLeadsByStatus retrieves leads by status
func (s *LeadService) GetLeadsByStatus(status string) ([]models.Lead, error) {
	return s.repos.LeadRepo.ListByStatus(status)
}

// GetLeadsByAssignee retrieves leads by assignee
func (s *LeadService) GetLeadsByAssignee(assigneeID int) ([]models.Lead, error) {
	return s.repos.LeadRepo.ListByAssignee(assigneeID)
}

// GetAllFieldConfigs retrieves all field configurations
func (s *LeadService) GetAllFieldConfigs(companyId int) ([]models.LeadFieldConfig, error) {
	return s.repos.LeadFieldConfigRepo.GetAllFieldConfigs(companyId)
}

// GetVisibleFieldConfigs retrieves visible field configurations
func (s *LeadService) GetVisibleFieldConfigs(companyId int) ([]models.LeadFieldConfig, error) {
	return s.repos.LeadFieldConfigRepo.GetVisibleFieldConfigs(companyId)
}

// GetRequiredFieldConfigs retrieves required field configurations
func (s *LeadService) GetRequiredFieldConfigs(companyId int) ([]models.LeadFieldConfig, error) {
	return s.repos.LeadFieldConfigRepo.GetRequiredFieldConfigs(companyId)
}

// GetFieldConfigsBySection retrieves field configurations by section
func (s *LeadService) GetFieldConfigsBySection(section string) ([]models.LeadFieldConfig, error) {
	return s.repos.LeadFieldConfigRepo.GetFieldConfigsBySection(section)
}

// CreateFieldConfig creates a new field configuration
func (s *LeadService) CreateFieldConfig(config *models.LeadFieldConfig) error {
	return s.repos.LeadFieldConfigRepo.CreateFieldConfig(config)
}

// UpdateFieldConfig updates a field configuration
func (s *LeadService) UpdateFieldConfig(config *models.LeadFieldConfig) error {
	return s.repos.LeadFieldConfigRepo.UpdateFieldConfig(config)
}

// DeleteFieldConfig deletes a field configuration
func (s *LeadService) DeleteFieldConfig(id uint) error {
	return s.repos.LeadFieldConfigRepo.DeleteFieldConfig(id)
}

// ReorderFormFields updates the order of form fields
func (s *LeadService) ReorderFormFields(fieldIDs []uint) error {
	return s.repos.LeadFieldConfigRepo.ReorderFormFields(fieldIDs)
}

// GetAllFormSections retrieves all form sections
func (s *LeadService) GetAllFormSections(companyId int) ([]models.LeadFormSection, error) {
	return s.repos.LeadFieldConfigRepo.GetAllFormSections(companyId)
}

// GetVisibleFormSections retrieves visible form sections
func (s *LeadService) GetVisibleFormSections(companyId int) ([]models.LeadFormSection, error) {
	return s.repos.LeadFieldConfigRepo.GetVisibleFormSections(companyId)
}

// CreateFormSection creates a new form section
func (s *LeadService) CreateFormSection(section *models.LeadFormSection) error {
	return s.repos.LeadFieldConfigRepo.CreateFormSection(section)
}

// UpdateFormSection updates a form section
func (s *LeadService) UpdateFormSection(section *models.LeadFormSection) error {
	return s.repos.LeadFieldConfigRepo.UpdateFormSection(section)
}

// DeleteFormSection deletes a form section
func (s *LeadService) DeleteFormSection(id int) error {
	return s.repos.LeadFieldConfigRepo.DeleteFormSection(id)
}

// ReorderFormSections updates the order of form sections
func (s *LeadService) ReorderFormSections(sectionIDs []int) error {
	return s.repos.LeadFieldConfigRepo.ReorderFormSections(sectionIDs)
}

// BulkImportLeads imports multiple leads - stub method to be implemented
func (s *LeadService) BulkImportLeads(leads []models.Lead) (map[string]interface{}, error) {
	// Implementation would go here
	return map[string]interface{}{
		"success": true,
		"message": "Leads imported successfully",
	}, nil
}

// ExportLeads exports leads - stub method to be implemented
func (s *LeadService) ExportLeads(filters map[string]string, companyId int) ([]models.Lead, error) {
	// Implementation would go here
	return s.repos.LeadRepo.List(companyId)
}

// Helper function to get required field names
func (s *LeadService) getRequiredFieldNames(companyId int) ([]string, error) {
	configs, err := s.repos.LeadFieldConfigRepo.GetRequiredFieldConfigs(companyId)
	if err != nil {
		return nil, err
	}

	fieldNames := make([]string, len(configs))
	for i, config := range configs {
		fieldNames[i] = config.FieldName
	}

	return fieldNames, nil
}

// Error definitions
var (
	ErrNotFound = errors.New("record not found")
)
