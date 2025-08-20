package repositories

import (
	"crm-app/backend/models"
	"errors"
	"fmt"
	"strconv"
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
// func (r *gormLeadRepository) List(companyId int) ([]models.Lead, error) {
// 	var leadss []models.Lead

// 	if err := r.db.Where("company_id = ?", companyId).Find(&leadss).Error; err != nil {
// 		return nil, err
// 	}

// 	var scoreTypes []models.ScoreType
// 	if error := r.db.Where("company_id = ?", companyId).Find(&scoreTypes).Error; error != nil {
// 		return nil, error
// 	}

// 	// Load related data for each lead
// 	for i := range leadss {
// 		// Load custom fields
// 		var customFields []models.LeadCustomField
// 		if err := r.db.Where("lead_id = ?", leadss[i].ID).Find(&customFields).Error; err != nil {
// 			return nil, err
// 		}
// 		leadss[i].CustomFields = customFields

// 		// Load tags
// 		var leadTags []models.LeadTag
// 		if err := r.db.Where("lead_id = ?", leadss[i].ID).Find(&leadTags).Error; err != nil {
// 			return nil, err
// 		}

// 		// Extract tag values
// 		tags := make([]string, len(leadTags))
// 		for j, tag := range leadTags {
// 			tags[j] = tag.Tag
// 		}
// 		leadss[i].Tags = tags

// 		matched := false

// 		if leadss[i].Score != nil {
// 			for _, st := range scoreTypes {
// 				if *leadss[i].Score >= st.MinScore && *leadss[i].Score <= st.MaxScore {
// 					leadss[i].Type = st.Type
// 					matched = true
// 					break
// 				}
// 			}
// 		}

// 		if !matched {
// 			leadss[i].Type = "cold"
// 		}
// 	}

// 	return leadss, nil
// }

// func (r *gormLeadRepository) List(companyId int) ([]models.GroupedLead, error) {
// 	type LeadFieldResult struct {
// 		SubmitID   uint   `gorm:"column:submit_id"`
// 		LeadID     uint   `gorm:"column:lead_id"`
// 		CrmFieldID uint   `gorm:"column:crm_field_id"`
// 		FieldName  string `gorm:"column:field_name"`
// 		FieldValue string `gorm:"column:field_value"`
// 	}
// 	var results []LeadFieldResult

// 	err := r.db.Table("leads").
// 		Select("crm_field_data.submit_id, leads.id as lead_id, crm_field_data.crm_field_id, lead_field_configs.field_name, crm_field_data.field_value").
// 		Joins("INNER JOIN crm_field_data ON crm_field_data.submit_id = leads.id").
// 		Joins("INNER JOIN lead_field_configs ON lead_field_configs.id = crm_field_data.crm_field_id").
// 		Scan(&results).Error

// 	if err != nil {
// 		return nil, err
// 	}

// 	var scoreTypes []models.ScoreType
// 	if error := r.db.Where("company_id = ?", companyId).Find(&scoreTypes).Error; error != nil {
// 		return nil, error
// 	}

// 	grouped := make(map[uint][]map[string]string)

// 	for _, r := range results {
// 		field := map[string]string{
// 			"fieldId":   fmt.Sprintf("%d", r.CrmFieldID),
// 			"fieldName": r.FieldName,
// 			"value":     r.FieldValue,
// 		}
// 		grouped[r.SubmitID] = append(grouped[r.SubmitID], field)
// 	}

// 	// Convert map to slice
// 	var finalResult []models.GroupedLead
// 	for submitID, fields := range grouped {
// 		finalResult = append(finalResult, models.GroupedLead{
// 			SubmitID: submitID,
// 			Fields:   fields,
// 		})
// 	}
// 	return finalResult, nil
// }

func (r *gormLeadRepository) List(companyId int) ([]models.GroupedLead, error) {

	var results []models.LeadFieldResult

	// Fetch lead fields
	err := r.db.Table("leads").
		Select("crm_field_data.submit_id, leads.id as lead_id, crm_field_data.crm_field_id, lead_field_configs.field_name, crm_field_data.field_value").
		Joins("INNER JOIN crm_field_data ON crm_field_data.submit_id = leads.id").
		Joins("INNER JOIN lead_field_configs ON lead_field_configs.id = crm_field_data.crm_field_id").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// Fetch scoring rules
	var scoreTypes []models.ScoreType
	if err := r.db.Where("company_id = ?", companyId).Find(&scoreTypes).Error; err != nil {
		return nil, err
	}

	// Group fields by submitId
	grouped := make(map[uint][]map[string]string)
	for _, r := range results {
		field := map[string]string{
			"fieldId":   fmt.Sprintf("%d", r.CrmFieldID),
			"fieldName": r.FieldName,
			"value":     r.FieldValue,
		}
		grouped[r.SubmitID] = append(grouped[r.SubmitID], field)
	}

	// Convert map to slice & apply scoring
	var finalResult []models.GroupedLead
	for submitID, fields := range grouped {
		leadScore := "cold" // default

		// try to find a numeric field for scoring
		for _, field := range fields {
			if scoreVal, err := strconv.Atoi(field["value"]); err == nil {
				for _, st := range scoreTypes {
					if scoreVal >= st.MinScore && scoreVal <= st.MaxScore {
						leadScore = st.Type
						break
					}
				}
				// stop after first valid scoring match
				break
			}
		}

		finalResult = append(finalResult, models.GroupedLead{
			SubmitID: submitID,
			Score:    leadScore,
			Fields:   fields,
		})
	}

	return finalResult, nil
}

// ListByStatus returns leads with the given status
func (r *gormLeadRepository) ListByStatus(status string) ([]models.GroupedLead, error) {
	var results []models.LeadFieldResult

	err := r.db.Table("leads").
		Select("crm_field_data.submit_id, leads.id as lead_id, crm_field_data.crm_field_id, lead_field_configs.field_name, crm_field_data.field_value").
		Joins("INNER JOIN crm_field_data ON crm_field_data.submit_id = leads.id").
		Joins("INNER JOIN lead_field_configs ON lead_field_configs.id = crm_field_data.crm_field_id").
		Where("leads.status = ?", status).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	grouped := make(map[uint][]map[string]string)

	for _, r := range results {
		field := map[string]string{
			"fieldId":   fmt.Sprintf("%d", r.CrmFieldID),
			"fieldName": r.FieldName,
			"value":     r.FieldValue,
		}
		grouped[r.SubmitID] = append(grouped[r.SubmitID], field)
	}

	// Convert map to slice
	var finalResult []models.GroupedLead
	for submitID, fields := range grouped {
		finalResult = append(finalResult, models.GroupedLead{
			SubmitID: submitID,
			Fields:   fields,
		})
	}
	return finalResult, nil
}

// ListByAssignee returns leads assigned to the given user
func (r *gormLeadRepository) ListByAssignee(assigneeID int) ([]models.GroupedLead, error) {

	var results []models.LeadFieldResult

	err := r.db.Table("leads").
		Select("crm_field_data.submit_id, leads.id as lead_id, crm_field_data.crm_field_id, lead_field_configs.field_name, crm_field_data.field_value").
		Joins("INNER JOIN crm_field_data ON crm_field_data.submit_id = leads.id").
		Joins("INNER JOIN lead_field_configs ON lead_field_configs.id = crm_field_data.crm_field_id").
		Where("leads.assigned_to_id = ?", assigneeID).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	grouped := make(map[uint][]map[string]string)

	for _, r := range results {
		field := map[string]string{
			"fieldId":   fmt.Sprintf("%d", r.CrmFieldID),
			"fieldName": r.FieldName,
			"value":     r.FieldValue,
		}
		grouped[r.SubmitID] = append(grouped[r.SubmitID], field)
	}

	// Convert map to slice
	var finalResult []models.GroupedLead
	for submitID, fields := range grouped {
		finalResult = append(finalResult, models.GroupedLead{
			SubmitID: submitID,
			Fields:   fields,
		})
	}
	return finalResult, nil
}

// Create creates a new lead
func (r *gormLeadRepository) Create(lead []models.CrmFieldData) error {
	return r.db.Create(&lead).Error
	// Start a transaction
	// tx := r.db.Begin()
	// if tx.Error != nil {
	// 	return tx.Error
	// }

	// Create the lead
	// if err := tx.Create(lead).Error; err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	// Save tags if provided
	// if len(lead.Tags) > 0 {
	// 	for _, tag := range lead.Tags {
	// 		leadTag := models.LeadTag{
	// 			LeadID:    lead.ID,
	// 			Tag:       tag,
	// 			CompanyId: lead.CompanyId,
	// 		}
	// 		if err := tx.Create(&leadTag).Error; err != nil {
	// 			tx.Rollback()
	// 			return err
	// 		}
	// 	}
	// }

	// Save custom fields if provided
	// if len(lead.CustomFields) > 0 {
	// 	for i := range lead.CustomFields {
	// 		lead.CustomFields[i].LeadID = lead.ID
	// 		lead.CustomFields[i].CompanyId = lead.CompanyId
	// 		if err := tx.Create(&lead.CustomFields[i]).Error; err != nil {
	// 			tx.Rollback()
	// 			return err
	// 		}
	// 	}
	// }

	// Commit the transaction
	// return tx.Commit().Error
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
	fmt.Println("missingFields 1", missingFields)
	fmt.Println("requiredFields 1", requiredFields)
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
	fmt.Println("missingFields 2", missingFields)
	if len(missingFields) > 0 {
		return fmt.Errorf("required fields missing: %s", strings.Join(missingFields, ", "))
	}

	return nil
}

func (r *gormLeadRepository) GetLastSubmitId() (int, error) {
	var lastId int
	err := r.db.Table("crm_field_data").
		Select("COALESCE(MAX(submit_id), 0)").
		Scan(&lastId).Error
	return lastId, err
}

func (r *gormLeadRepository) CreateMainLead(lead *models.Lead) error {
	fmt.Println("inline")
	return r.db.Create(lead).Error
}
