package repositories

import (
	"crm-app/backend/models"

	"gorm.io/gorm"
)

// GormLeadFieldConfigRepository implements LeadFieldConfigRepository with GORM
// type GormLeadFieldConfigRepository struct {
// 	db *gorm.DB
// }

// NewLeadFieldConfigRepository creates a new lead field config repository
// func NewLeadFieldConfigRepository(db *gorm.DB) models.LeadFieldConfigRepository {
// 	return &GormLeadFieldConfigRepository{db: db}
// }

// GetAllFieldConfigs retrieves all field configurations
func (r *GormLeadFieldConfigRepository) GetAllFieldConfigs(companyId int) ([]models.LeadFieldConfig, error) {
	var configs []models.LeadFieldConfig
	if err := r.db.Where("company_id=?", companyId).Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// GetVisibleFieldConfigs retrieves visible field configurations
func (r *GormLeadFieldConfigRepository) GetVisibleFieldConfigs(companyId int) ([]models.LeadFieldConfig, error) {
	var configs []models.LeadFieldConfig
	if err := r.db.Where("visible = ? AND company_id= ?", true, companyId).Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// GetRequiredFieldConfigs retrieves required field configurations
func (r *GormLeadFieldConfigRepository) GetRequiredFieldConfigs(companyId int) ([]models.LeadFieldConfig, error) {
	var configs []models.LeadFieldConfig
	if err := r.db.Where("required = ? AND company_id= ?", true, companyId).Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// GetFieldConfigsBySection retrieves field configurations by section
func (r *GormLeadFieldConfigRepository) GetFieldConfigsBySection(section string) ([]models.LeadFieldConfig, error) {
	var configs []models.LeadFieldConfig
	if err := r.db.Where("section = ?", section).Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// GetFieldConfig retrieves a field configuration by ID
func (r *GormLeadFieldConfigRepository) GetFieldConfig(id int) (*models.LeadFieldConfig, error) {
	var config models.LeadFieldConfig
	if err := r.db.First(&config, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &config, nil
}

// CreateFieldConfig creates a new field configuration
func (r *GormLeadFieldConfigRepository) CreateFieldConfig(config *models.LeadFieldConfig) error {
	return r.db.Create(config).Error
}

// UpdateFieldConfig updates a field configuration
func (r *GormLeadFieldConfigRepository) UpdateFieldConfig(config *models.LeadFieldConfig) error {
	return r.db.Save(config).Error
}

// DeleteFieldConfig deletes a field configuration
func (r *GormLeadFieldConfigRepository) DeleteFieldConfig(id uint) error {
	return r.db.Delete(&models.LeadFieldConfig{}, id).Error
}

// ReorderFormFields updates the order of form fields
func (r *GormLeadFieldConfigRepository) ReorderFormFields(fieldIDs []uint) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for i, id := range fieldIDs {
		if err := tx.Model(&models.LeadFieldConfig{}).Where("id = ?", id).Update("order_index", i).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// GetAllFormSections retrieves all form sections
func (r *GormLeadFieldConfigRepository) GetAllFormSections(companyId int) ([]models.LeadFormSection, error) {
	var sections []models.LeadFormSection
	if err := r.db.Where("company_id = ?", companyId).Find(&sections).Error; err != nil {
		return nil, err
	}
	return sections, nil
}

// GetVisibleFormSections retrieves visible form sections
func (r *GormLeadFieldConfigRepository) GetVisibleFormSections(companyId int) ([]models.LeadFormSection, error) {
	var sections []models.LeadFormSection
	if err := r.db.Where("visible = ? AND company_id = ?", true, companyId).Find(&sections).Error; err != nil {
		return nil, err
	}
	return sections, nil
}

// CreateFormSection creates a new form section
func (r *GormLeadFieldConfigRepository) CreateFormSection(section *models.LeadFormSection) error {
	return r.db.Create(section).Error
}

// UpdateFormSection updates a form section
func (r *GormLeadFieldConfigRepository) UpdateFormSection(section *models.LeadFormSection) error {
	return r.db.Save(section).Error
}

// DeleteFormSection deletes a form section
func (r *GormLeadFieldConfigRepository) DeleteFormSection(id int) error {
	return r.db.Delete(&models.LeadFormSection{}, id).Error
}

// ReorderFormSections updates the order of form sections
func (r *GormLeadFieldConfigRepository) ReorderFormSections(sectionIDs []int) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for i, id := range sectionIDs {
		if err := tx.Model(&models.LeadFormSection{}).Where("id = ?", id).Update("order_index", i).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *GormScoreRepository) ScoreUpdateRepo(config []models.ScoreType) error {
	for _, config := range config {
		if err := r.DB.Save(&config).Error; err != nil {
			return err // or collect all errors if you want to return multiple
		}
	}
	return nil
}

// GetFormStructure retrieves the form structure
// func (r *GormLeadFieldConfigRepository) GetFormStructure() (map[string]interface{}, error) {
// 	sections, err := r.GetVisibleFormSections()
// 	if err != nil {
// 		return nil, err
// 	}

// 	result := map[string]interface{}{
// 		"sections": sections,
// 	}

// 	// Get fields for each section
// 	for i, section := range sections {
// 		fields, err := r.GetFieldConfigsBySection(section.Name)
// 		if err != nil {
// 			return nil, err
// 		}
// 		sectionMap := result["sections"].([]models.LeadFormSection)
// 		// sectionMap[i].Fields = fields
// 	}

// 	return result, nil
// }

func (r *GormLeadFieldConfigRepository) GetFormStructure(companyId int) (map[string]interface{}, error) {
	// Get all sections
	sections, err := r.GetVisibleFormSections(companyId)
	if err != nil {
		return nil, err
	}

	// Get all fields
	fields, err := r.GetVisibleFieldConfigs(companyId)
	if err != nil {
		return nil, err
	}

	// Group fields by section
	fieldsBySection := make(map[string][]models.LeadFieldConfig)
	for _, field := range fields {
		if _, ok := fieldsBySection[field.Section]; !ok {
			fieldsBySection[field.Section] = []models.LeadFieldConfig{}
		}
		fieldsBySection[field.Section] = append(fieldsBySection[field.Section], field)
	}

	// Build the structure
	result := make(map[string]interface{})
	sectionsList := []map[string]interface{}{}

	// Default section for fields without a section
	if fields, ok := fieldsBySection["default"]; ok {
		defaultSection := map[string]interface{}{
			"name":   "default",
			"label":  "Default",
			"fields": fields,
		}
		sectionsList = append(sectionsList, defaultSection)
	}

	// Add all defined sections
	for _, section := range sections {
		sectionData := map[string]interface{}{
			"id":          section.ID,
			"name":        section.Name,
			"label":       section.Label,
			"description": section.Description,
			"orderIndex":  section.OrderIndex,
			"fields":      fieldsBySection[section.Name],
		}
		sectionsList = append(sectionsList, sectionData)
	}

	result["sections"] = sectionsList
	return result, nil
}

// func (r *GormScoreRepository) ScoreUpdateRepo(config []models.ScoreType) error {
// 	for _, config := range config {
// 		if err := r.DB.Save(&config).Error; err != nil {
// 			return err // or collect all errors if you want to return multiple
// 		}
// 	}
// 	return nil
// }
