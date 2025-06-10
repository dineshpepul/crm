package repositories

import (
	"crm-app/backend/models"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// gormLeadRepository implements LeadRepository with GORM
// type gormLeadRepository struct {
// 	db *gorm.DB
// }

// NewLeadRepository creates a new lead repository
// func NewLeadRepository(db *gorm.DB) models.LeadRepository {
// 	return &gormLeadRepository{db: db}
// }

// FindByID finds a lead by ID
func (r *gormLeadRepository) FindByID(id int) (*models.Lead, error) {
	var lead models.Lead
	result := r.db.First(&lead, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	// Load custom fields
	var customFields []models.LeadCustomField
	if err := r.db.Where("lead_id = ?", lead.ID).Find(&customFields).Error; err != nil {
		return nil, err
	}
	lead.CustomFields = customFields

	// Load tags
	var leadTags []models.LeadTag
	if err := r.db.Where("lead_id = ?", lead.ID).Find(&leadTags).Error; err != nil {
		return nil, err
	}

	// Extract tag values
	tags := make([]string, len(leadTags))
	for i, tag := range leadTags {
		tags[i] = tag.Tag
	}
	lead.Tags = tags

	return &lead, nil
}

// List returns all leads
func (r *gormLeadRepository) List(companyId int) ([]models.Lead, error) {
	var leads []models.Lead
	if err := r.db.Find(&leads).Error; err != nil {
		return nil, err
	}

	// Load related data for each lead
	for i := range leads {
		// Load custom fields
		var customFields []models.LeadCustomField
		if err := r.db.Where("lead_id = ?", leads[i].ID).Find(&customFields).Error; err != nil {
			return nil, err
		}
		leads[i].CustomFields = customFields

		// Load tags
		var leadTags []models.LeadTag
		if err := r.db.Where("lead_id = ?", leads[i].ID).Find(&leadTags).Error; err != nil {
			return nil, err
		}

		// Extract tag values
		tags := make([]string, len(leadTags))
		for j, tag := range leadTags {
			tags[j] = tag.Tag
		}
		leads[i].Tags = tags
	}

	return leads, nil
}

// ListByStatus returns leads with the given status
func (r *gormLeadRepository) ListByStatus(status string) ([]models.Lead, error) {
	var leads []models.Lead
	if err := r.db.Where("status = ?", status).Find(&leads).Error; err != nil {
		return nil, err
	}

	var scoreTypes []models.ScoreType
	if error := r.db.Find(&scoreTypes).Error; error != nil {
		return nil, error
	}

	// Load related data for each lead (custom fields and tags)
	for i := range leads {
		var customFields []models.LeadCustomField
		if err := r.db.Where("lead_id = ?", leads[i].ID).Find(&customFields).Error; err != nil {
			return nil, err
		}
		leads[i].CustomFields = customFields

		matched := false

		if leads[i].Score != nil {
			for _, st := range scoreTypes {
				if *leads[i].Score >= st.MinScore && *leads[i].Score <= st.MaxScore {
					leads[i].Type = st.Type
					matched = true
					break
				}
			}
		}

		if !matched {
			leads[i].Type = "cold"
		}

		var leadTags []models.LeadTag
		if err := r.db.Where("lead_id = ?", leads[i].ID).Find(&leadTags).Error; err != nil {
			return nil, err
		}

		tags := make([]string, len(leadTags))
		for j, tag := range leadTags {
			tags[j] = tag.Tag
		}
		leads[i].Tags = tags
	}

	return leads, nil
}

// ListByAssignee returns leads assigned to the given user
func (r *gormLeadRepository) ListByAssignee(assigneeID int) ([]models.Lead, error) {
	var leads []models.Lead
	if err := r.db.Where("assigned_to_id = ?", assigneeID).Find(&leads).Error; err != nil {
		return nil, err
	}

	// Load related data for each lead (custom fields and tags)
	for i := range leads {
		var customFields []models.LeadCustomField
		if err := r.db.Where("lead_id = ?", leads[i].ID).Find(&customFields).Error; err != nil {
			return nil, err
		}
		leads[i].CustomFields = customFields

		var leadTags []models.LeadTag
		if err := r.db.Where("lead_id = ?", leads[i].ID).Find(&leadTags).Error; err != nil {
			return nil, err
		}

		tags := make([]string, len(leadTags))
		for j, tag := range leadTags {
			tags[j] = tag.Tag
		}
		leads[i].Tags = tags
	}

	return leads, nil
}

// Create creates a new lead
func (r *gormLeadRepository) Create(lead *models.Lead) error {
	// Start a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Ensure default status
	if lead.Status == "" {
		lead.Status = "new"
	}

	// Create the lead
	if err := tx.Create(lead).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Save tags if provided
	if len(lead.Tags) > 0 {
		for _, tag := range lead.Tags {
			leadTag := models.LeadTag{
				LeadID:    lead.ID,
				Tag:       tag,
				CompanyId: lead.CompanyId,
			}
			if err := tx.Create(&leadTag).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Save custom fields if provided
	if len(lead.CustomFields) > 0 {
		for i := range lead.CustomFields {
			lead.CustomFields[i].LeadID = lead.ID
			lead.CustomFields[i].CompanyId = lead.CompanyId
			if err := tx.Create(&lead.CustomFields[i]).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Commit the transaction
	return tx.Commit().Error
}

// Update updates a lead
func (r *gormLeadRepository) Update(lead *models.Lead) error {
	// Start a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Update the lead
	if err := tx.Save(lead).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete existing tags
	if err := tx.Where("lead_id = ?", lead.ID).Delete(&models.LeadTag{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Save new tags if provided
	if len(lead.Tags) > 0 {
		for _, tag := range lead.Tags {
			leadTag := models.LeadTag{
				LeadID: lead.ID,
				Tag:    tag,
			}
			if err := tx.Create(&leadTag).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Delete existing custom fields
	if err := tx.Where("lead_id = ?", lead.ID).Delete(&models.LeadCustomField{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Save new custom fields if provided
	if len(lead.CustomFields) > 0 {
		for i := range lead.CustomFields {
			lead.CustomFields[i].LeadID = lead.ID
			if err := tx.Create(&lead.CustomFields[i]).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Commit the transaction
	return tx.Commit().Error
}

// Delete deletes a lead
func (r *gormLeadRepository) Delete(id int) error {
	return r.db.Delete(&models.Lead{}, id).Error
}

// ValidateLeadFields validates that all required fields are present in the lead
func (r *gormLeadRepository) ValidateLeadFields(lead *models.Lead, requiredFields []string) error {
	missingFields := []string{}

	for _, field := range requiredFields {
		switch field {
		case "name":
			if lead.Name == "" {
				missingFields = append(missingFields, "name")
			}
		case "email":
			if lead.Email == "" {
				missingFields = append(missingFields, "email")
			}
		case "phone":
			if lead.Phone == "" {
				missingFields = append(missingFields, "phone")
			}
		case "company":
			if lead.Company == "" {
				missingFields = append(missingFields, "company")
			}
		case "source":
			if lead.Source == "" {
				missingFields = append(missingFields, "source")
			}
		case "status":
			if lead.Status == "" {
				missingFields = append(missingFields, "status")
			}
		default:
			// Check custom fields
			found := false
			for _, cf := range lead.CustomFields {
				if cf.FieldName == field && cf.FieldValue != "" {
					found = true
					break
				}
			}
			if !found {
				missingFields = append(missingFields, field)
			}
		}
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("required fields missing: %s", strings.Join(missingFields, ", "))
	}

	return nil
}
