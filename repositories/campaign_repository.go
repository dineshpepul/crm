package repositories

import (
	"crm-app/backend/models"

	"gorm.io/gorm"
)

// type gormCampaignRepository struct {
// 	db *gorm.DB
// }

// // NewCampaignRepository creates a new campaign repository
// func NewCampaignRepository(db *gorm.DB) models.CampaignRepository {
// 	return &gormCampaignRepository{db: db}
// }

// GetCampaigns returns campaigns
func (r *gormCampaignRepository) GetCampaigns(offset int, limit int) ([]models.Campaign, error) {
	var campaigns []models.Campaign
	err := r.db.Offset(offset).Limit(limit).Find(&campaigns).Error
	return campaigns, err
}

// GetCampaignByID returns a campaign by ID
func (r *gormCampaignRepository) GetCampaignByID(id int) (*models.Campaign, error) {
	var campaign models.Campaign
	err := r.db.First(&campaign, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &campaign, nil
}

// CreateCampaign creates a new campaign
func (r *gormCampaignRepository) CreateCampaign(campaign *models.Campaign) error {
	return r.db.Create(campaign).Error
}

// UpdateCampaign updates a campaign
func (r *gormCampaignRepository) UpdateCampaign(campaign *models.Campaign) error {
	return r.db.Save(campaign).Error
}

// DeleteCampaign deletes a campaign
func (r *gormCampaignRepository) DeleteCampaign(id int) error {
	return r.db.Delete(&models.Campaign{}, id).Error
}

// GetCampaignStats returns campaign statistics
func (r *gormCampaignRepository) GetCampaignStats(id int) (map[string]interface{}, error) {
	// This is a stub implementation
	stats := map[string]interface{}{
		"sent":      0,
		"opened":    0,
		"clicked":   0,
		"bounced":   0,
		"leads":     0,
		"revenue":   0,
		"roi":       0,
		"ctr":       0,
		"open_rate": 0,
	}

	return stats, nil
}

// GetLeadsForCampaign returns leads for a campaign
func (r *gormCampaignRepository) GetLeadsForCampaign(id int) ([]models.Lead, error) {
	var leads []models.Lead
	err := r.db.Table("campaign_leads").
		Select("leads.*").
		Joins("JOIN leads ON campaign_leads.lead_id = leads.id").
		Where("campaign_leads.campaign_id = ?", id).
		Find(&leads).Error
	return leads, err
}

// AssignLeadsToCampaign assigns leads to a campaign
func (r *gormCampaignRepository) AssignLeadsToCampaign(campaignID int, leadIDs []int) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, leadID := range leadIDs {
		if err := tx.Exec("INSERT INTO campaign_leads (campaign_id, lead_id) VALUES (?, ?)", campaignID, leadID).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// RemoveLeadsFromCampaign removes leads from a campaign
func (r *gormCampaignRepository) RemoveLeadsFromCampaign(campaignID int, leadIDs []int) error {
	return r.db.Where("campaign_id = ? AND lead_id IN ?", campaignID, leadIDs).Delete("campaign_leads").Error
}

// GetTemplates returns campaign templates
func (r *gormCampaignRepository) GetTemplates(offset int, limit int) ([]models.CampaignTemplate, error) {
	var templates []models.CampaignTemplate
	err := r.db.Offset(offset).Limit(limit).Find(&templates).Error
	return templates, err
}

// GetTemplateByID returns a campaign template by ID
func (r *gormCampaignRepository) GetTemplateByID(id int) (*models.CampaignTemplate, error) {
	var template models.CampaignTemplate
	err := r.db.First(&template, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &template, nil
}

// CreateTemplate creates a new campaign template
func (r *gormCampaignRepository) CreateTemplate(template *models.CampaignTemplate) error {
	return r.db.Create(template).Error
}

// UpdateTemplate updates a campaign template
func (r *gormCampaignRepository) UpdateTemplate(template *models.CampaignTemplate) error {
	return r.db.Save(template).Error
}

// DeleteTemplate deletes a campaign template
func (r *gormCampaignRepository) DeleteTemplate(id int) error {
	return r.db.Delete(&models.CampaignTemplate{}, id).Error
}
