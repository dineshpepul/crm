package repositories

import (
	"gorm.io/gorm"
)

// GormLeadFieldConfigRepository implements LeadFieldConfigRepository with GORM
type GormLeadFieldConfigRepository struct {
	db *gorm.DB
}

// NewLeadFieldConfigRepository creates a new lead field config repository
// func NewLeadFieldConfigRepository(db *gorm.DB) models.LeadFieldConfigRepository {
// 	return &GormLeadFieldConfigRepository{db: db}
// }

// GetAllFieldConfigs returns all field configurations
// func (r *GormLeadFieldConfigRepository) GetAllFieldConfigs() ([]models.LeadFieldConfig, error) {
// 	var configs []models.LeadFieldConfig
// 	err := r.db.Order("section, order_index").Find(&configs).Error
// 	return configs, err
// }

// // GetVisibleFieldConfigs returns visible field configurations
// func (r *GormLeadFieldConfigRepository) GetVisibleFieldConfigs() ([]models.LeadFieldConfig, error) {
// 	var configs []models.LeadFieldConfig
// 	err := r.db.Where("visible = ?", true).Order("section, order_index").Find(&configs).Error
// 	return configs, err
// }

// // GetRequiredFieldConfigs returns required field configurations
// func (r *GormLeadFieldConfigRepository) GetRequiredFieldConfigs() ([]models.LeadFieldConfig, error) {
// 	var configs []models.LeadFieldConfig
// 	err := r.db.Where("required = ?", true).Order("section, order_index").Find(&configs).Error
// 	return configs, err
// }

// // GetFieldConfigsBySection returns field configurations by section
// func (r *GormLeadFieldConfigRepository) GetFieldConfigsBySection(section string) ([]models.LeadFieldConfig, error) {
// 	var configs []models.LeadFieldConfig
// 	err := r.db.Where("section = ?", section).Order("order_index").Find(&configs).Error
// 	return configs, err
// }

// // CreateFieldConfig creates a new field configuration
// func (r *GormLeadFieldConfigRepository) CreateFieldConfig(config *models.LeadFieldConfig) error {
// 	return r.db.Create(config).Error
// }

// // UpdateFieldConfig updates a field configuration
// func (r *GormLeadFieldConfigRepository) UpdateFieldConfig(config *models.LeadFieldConfig) error {
// 	return r.db.Save(config).Error
// }

// // DeleteFieldConfig deletes a field configuration
// func (r *GormLeadFieldConfigRepository) DeleteFieldConfig(id uint) error {
// 	return r.db.Delete(&models.LeadFieldConfig{}, id).Error
// }

// // ReorderFormFields updates the order of form fields
// func (r *GormLeadFieldConfigRepository) ReorderFormFields(fieldIDs []uint) error {
// 	tx := r.db.Begin()
// 	if tx.Error != nil {
// 		return tx.Error
// 	}

// 	for i, id := range fieldIDs {
// 		if err := tx.Model(&models.LeadFieldConfig{}).Where("id = ?", id).Update("order_index", i).Error; err != nil {
// 			tx.Rollback()
// 			return err
// 		}
// 	}

// 	return tx.Commit().Error
// }
