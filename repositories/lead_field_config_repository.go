package repositories

import (
	"crm-app/backend/models"
	"errors"
	"fmt"
	"time"

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
	return r.db.Omit("CreatedAt").Save(config).Error
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
	// If no sections found, call insert function
	if len(sections) == 0 {
		fmt.Println("No sections found, inserting default sections...")
		if err := r.db.Where("company_id = ?", companyId).Find(&sections).Error; err != nil {
			return nil, err
		}
		if len(sections) == 0 {
			err := r.InsertDefaultFormSections(companyId)
			if err != nil {
				return nil, err
			}
		}
		// Optionally re-fetch the inserted sections
		if err := r.db.Where("visible = ? AND company_id = ?", true, companyId).Find(&sections).Error; err != nil {
			return nil, err
		}
	}

	return sections, nil
}

func (r *GormLeadFieldConfigRepository) InsertDefaultFormSections(companyId int) error {
	defaultStages := []struct {
		Label     string
		Title     string
		FormField []models.LeadFieldConfig
	}{
		{
			Label: "Lead Information",
			Title: "lead_information",
			FormField: []models.LeadFieldConfig{
				{FieldName: "name", DisplayName: "Name", FieldType: "text", Options: `[""]`, Required: true, Visible: true, Section: "lead_information", OrderIndex: 1, Placeholder: "Enter your name"},
				{FieldName: "title", DisplayName: "Title", FieldType: "text", Options: `[""]`, Required: true, Visible: true, Section: "lead_information", OrderIndex: 1, Placeholder: "Enter lead title"},
				{FieldName: "Lead_source", DisplayName: "Lead Source", FieldType: "select", Options: `["Advertisement", "Cold Call", "Employee Referral", "External Referral", "Online Store", "X (Twitter)", "LinkedIn", "Facebook", "Instagram", "Whatsapp", "Partner", "Public Relations", "Sales Email Alias", "Seminar Partner", "Internal Seminar", "Trade Show", "Web Download", "Web Research", "Chat", "Website Form", "Cold Email"]`, Required: false, Visible: true, Section: "lead_information", OrderIndex: 1, Placeholder: "Choose the lead source"},
				{FieldName: "lead_status", DisplayName: "Lead Status", FieldType: "select", Options: `["Attempted to contact", "Contact in future", "Contacted", "Junk Lead", "Lost Lead", "Not Contacted", "Pre-Qualified", "Not Qualified", "Hot"]`, Required: false, Visible: true, Section: "lead_information", OrderIndex: 1, Placeholder: "Choose the lead status"},
				{FieldName: "industry", DisplayName: "Industry", FieldType: "select", Options: `["ASP (Application Service Provider)", "Data/Telecom OEM", "ERP (Enterprise Resource Planning)", "Government/Military", "Large Enterprise", "ManagementISV", "MSP (Management Service Provider)", "Network Equipment Enterprise", "Non-management ISV", "Optical Networking", "Service Provider", "Small/Medium Enterprise", "Storage Equipment", "Storage Service Provider", "Systems Integrator", "Wireless Industry", "ERP", "Management ISV"]`, Required: false, Visible: true, Section: "lead_information", OrderIndex: 1, Placeholder: "Choose the lead industry"},
				{FieldName: "company", DisplayName: "Company", FieldType: "text", Options: `[""]`, Required: false, Visible: true, Section: "lead_information", OrderIndex: 1, Placeholder: "Enter company name"},
				{FieldName: "no_of_employee", DisplayName: "No. of employee", FieldType: "text", Options: `[""]`, Required: false, Visible: true, Section: "lead_information", OrderIndex: 1, Placeholder: "Enter no. of employee"},
				{FieldName: "email", DisplayName: "Email", FieldType: "email", Options: `[""]`, Required: true, Visible: false, Section: "lead_information", OrderIndex: 1, Placeholder: "Enter email"},
				{FieldName: "mobile_no_1", DisplayName: "Mobile no 1", FieldType: "phone", Options: `[""]`, Required: false, Visible: true, Section: "lead_information", OrderIndex: 1, Placeholder: "Enter your mobile no"},
				{FieldName: "mobile_no_2", DisplayName: "Mobile no 2", FieldType: "phone", Options: `[""]`, Required: false, Visible: true, Section: "lead_information", OrderIndex: 1, Placeholder: "Enter your mobile no"},
				{FieldName: "gender", DisplayName: "Gender", FieldType: "text", Options: `[""]`, Required: false, Visible: true, Section: "lead_information", OrderIndex: 1, Placeholder: "Enter gender"},
				{FieldName: "pan_card", DisplayName: "Pan Card", FieldType: "text", Options: `[""]`, Required: false, Visible: true, Section: "lead_information", OrderIndex: 1, Placeholder: "Enter pan number"},
			},
		},
		{
			Label: "Address Information",
			Title: "address_information",
			FormField: []models.LeadFieldConfig{
				{FieldName: "street", DisplayName: "Street", FieldType: "text", Placeholder: "Enter your street", Options: `[""]`, Required: false, Visible: true, Section: "address_information", OrderIndex: 1},
				{FieldName: "city", DisplayName: "City", FieldType: "text", Placeholder: "Enter your city", Options: `[""]`, Required: false, Visible: true, Section: "address_information", OrderIndex: 2},
				{FieldName: "state", DisplayName: "State", FieldType: "text", Placeholder: "Enter your state", Options: `[""]`, Required: false, Visible: true, Section: "address_information", OrderIndex: 3},
				{FieldName: "zip_code", DisplayName: "Zip Code", FieldType: "text", Placeholder: "Enter your zip code", Options: `[""]`, Required: false, Visible: true, Section: "address_information", OrderIndex: 4},
				{FieldName: "country", DisplayName: "Country", FieldType: "text", Placeholder: "Enter your country", Options: `[""]`, Required: false, Visible: true, Section: "address_information", OrderIndex: 5},
				{FieldName: "website_url", DisplayName: "Website Url", FieldType: "text", Placeholder: "Enter your website url", Options: `[""]`, Required: false, Visible: true, Section: "address_information", OrderIndex: 6},
			},
		},
		{
			Label: "Description Information",
			Title: "description_information",
			FormField: []models.LeadFieldConfig{
				{FieldName: "description_information", DisplayName: "Description Information", FieldType: "text", Placeholder: "Enter your description", Options: `[""]`, Required: false, Visible: true, Section: "description_information", OrderIndex: 1},
			},
		},
	}
	var scoreTypes []models.ScoreType
	if err := r.db.Where("company_id = ?", companyId).Find(&scoreTypes).Error; err != nil {
		return nil
	}

	if len(scoreTypes) == 0 {
		defaultTypes := []models.ScoreType{
			{Type: "cold", MinScore: 0, MaxScore: 30, CompanyId: companyId},
			{Type: "warm", MinScore: 31, MaxScore: 60, CompanyId: companyId},
			{Type: "hot", MinScore: 61, MaxScore: 100, CompanyId: companyId},
		}

		if err := r.db.Create(&defaultTypes).Error; err != nil {
			return nil
		}

		scoreTypes = defaultTypes
	}

	now := time.Now()
	for i, stage := range defaultStages {
		stageData := models.LeadFormSection{
			CompanyId:  companyId,
			Name:       stage.Title,
			Label:      stage.Label,
			OrderIndex: i + 1,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		if err := r.db.Create(&stageData).Error; err != nil {
			return err
		}
		fmt.Println("id", stageData.ID)
		for _, field := range stage.FormField {
			fieldData := models.LeadFieldConfig{
				CompanyId:   companyId,
				FieldName:   field.FieldName,
				DisplayName: field.DisplayName,
				FieldType:   field.FieldType,
				Options:     field.Options,
				Required:    field.Required,
				Visible:     field.Visible,
				Section:     stage.Title,
				SectionId:   int(stageData.ID),
				OrderIndex:  field.OrderIndex,
				Placeholder: field.Placeholder,
				CreatedAt:   now,
				UpdatedAt:   now,
			}

			if err := r.db.Create(&fieldData).Error; err != nil {
				return err
			}
		}
	}
	return nil
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
	var LeadFormSection models.LeadFormSection
	if err := r.db.Model(&LeadFormSection).Select("name,company_id").Where("id = ?", id).First(&LeadFormSection).Error; err != nil {
		return nil
	}
	fmt.Println("LeadFormSection", LeadFormSection)
	if error := r.db.Where("company_id = ? AND section = ?", LeadFormSection.CompanyId, LeadFormSection.Name).Delete(&models.LeadFieldConfig{}).Error; error != nil {
		return nil
	}

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
	// for _, config := range config {
	// 	if err := r.DB.Save(&config).Error; err != nil {
	// 		return err // or collect all errors if you want to return multiple
	// 	}
	// }
	// return nil
	for _, cfg := range config {
		var existing models.ScoreType

		// Check if score type exists for company_id + type
		err := r.DB.Where("company_id = ? AND type = ?", cfg.CompanyId, cfg.Type).
			First(&existing).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Not found → create new
				if err := r.DB.Create(&cfg).Error; err != nil {
					return err
				}
			} else {
				return err // real DB error
			}
		} else {
			// Found → update existing
			existing.MinScore = cfg.MinScore
			existing.MaxScore = cfg.MaxScore
			existing.Type = cfg.Type
			if err := r.DB.Save(&existing).Error; err != nil {
				return err
			}
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
