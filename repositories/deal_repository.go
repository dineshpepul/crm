package repositories

import (
	"crm-app/backend/models"
	"errors"

	"gorm.io/gorm"
)

// gormDealRepository implements DealRepository with GORM
// type gormDealRepository struct {
// 	db *gorm.DB
// }

// NewDealRepository creates a new deal repository
// func NewDealRepository(db *gorm.DB) models.DealRepository {
// 	return &gormDealRepository{db: db}
// }

// FindByID finds a deal by ID
func (r *gormDealRepository) FindByID(id int) (*models.Deal, error) {
	var deal models.Deal
	result := r.db.First(&deal, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &deal, nil
}

// List returns deals with pagination and filters
func (r *gormDealRepository) List(offset int, limit int, filters map[string]interface{}) ([]models.Deal, error) {
	var deals []models.Deal
	query := r.db

	// Apply filters if provided
	if filters != nil {
		for key, value := range filters {
			query = query.Where(key+" = ?", value)
		}
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&deals).Error; err != nil {
		return nil, err
	}
	return deals, nil
}

// FindByLead returns deals for a specific lead
func (r *gormDealRepository) FindByLead(leadID int) ([]models.Deal, error) {
	var deals []models.Deal
	if err := r.db.Where("lead_id = ?", leadID).Find(&deals).Error; err != nil {
		return nil, err
	}
	return deals, nil
}

// Create creates a new deal
func (r *gormDealRepository) Create(deal *models.Deal) error {
	return r.db.Create(deal).Error
}

// Update updates an existing deal
func (r *gormDealRepository) Update(deal *models.Deal) error {
	return r.db.Save(deal).Error
}

// Delete deletes a deal
func (r *gormDealRepository) Delete(id int) error {
	return r.db.Delete(&models.Deal{}, id).Error
}

// GetDealPipeline returns the deal pipeline statistics
func (r *gormDealRepository) GetDealPipeline() ([]map[string]interface{}, error) {
	type PipelineStage struct {
		Stage      string  `json:"stage"`
		Count      int     `json:"count"`
		TotalValue float64 `json:"total_value"`
	}

	var results []PipelineStage

	if err := r.db.Model(&models.Deal{}).
		Select("stage, COUNT(*) as count, SUM(amount) as total_value").
		Group("stage").
		Order("FIELD(stage, 'lead', 'qualified', 'proposal', 'negotiation', 'closed_won', 'closed_lost')").
		Find(&results).Error; err != nil {
		return nil, err
	}

	// Convert to map[string]interface{} for flexibility
	pipeline := make([]map[string]interface{}, len(results))
	for i, stage := range results {
		pipeline[i] = map[string]interface{}{
			"stage":       stage.Stage,
			"count":       stage.Count,
			"total_value": stage.TotalValue,
		}
	}

	return pipeline, nil
}
