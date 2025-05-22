package repositories

import (
	"crm-app/backend/models"

	"gorm.io/gorm"
)

// type gormNurtureRepository struct {
// 	db *gorm.DB
// }

// NewNurtureRepository creates a new nurture repository
// func NewNurtureRepository(db *gorm.DB) models.NurtureRepository {
// 	return &gormNurtureRepository{db: db}
// }

// GetSequences returns nurture sequences
func (r *gormNurtureRepository) GetSequences(offset int, limit int) ([]models.NurtureSequence, error) {
	var sequences []models.NurtureSequence
	err := r.db.Offset(offset).Limit(limit).Find(&sequences).Error
	return sequences, err
}

// GetSequenceByID returns a nurture sequence by ID
func (r *gormNurtureRepository) GetSequenceByID(id int) (*models.NurtureSequence, error) {
	var sequence models.NurtureSequence
	err := r.db.First(&sequence, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &sequence, nil
}

// CreateSequence creates a new nurture sequence
func (r *gormNurtureRepository) CreateSequence(sequence *models.NurtureSequence) error {
	return r.db.Create(sequence).Error
}

// UpdateSequence updates a nurture sequence
func (r *gormNurtureRepository) UpdateSequence(sequence *models.NurtureSequence) error {
	return r.db.Save(sequence).Error
}

// DeleteSequence deletes a nurture sequence
func (r *gormNurtureRepository) DeleteSequence(id int) error {
	return r.db.Delete(&models.NurtureSequence{}, id).Error
}

// GetStepsBySequence returns steps by sequence ID
func (r *gormNurtureRepository) GetStepsBySequence(sequenceID int) ([]models.NurtureStep, error) {
	var steps []models.NurtureStep
	err := r.db.Where("sequence_id = ?", sequenceID).Order("order_index").Find(&steps).Error
	return steps, err
}

// CreateStep creates a new nurture step
func (r *gormNurtureRepository) CreateStep(step *models.NurtureStep) error {
	return r.db.Create(step).Error
}

// UpdateStep updates a nurture step
func (r *gormNurtureRepository) UpdateStep(step *models.NurtureStep) error {
	return r.db.Save(step).Error
}

// DeleteStep deletes a nurture step
func (r *gormNurtureRepository) DeleteStep(id int) error {
	return r.db.Delete(&models.NurtureStep{}, id).Error
}

// GetEnrollments returns enrollments by sequence ID
func (r *gormNurtureRepository) GetEnrollments(sequenceID int, offset int, limit int) ([]models.NurtureEnrollment, error) {
	var enrollments []models.NurtureEnrollment
	err := r.db.Where("sequence_id = ?", sequenceID).Offset(offset).Limit(limit).Find(&enrollments).Error
	return enrollments, err
}

// EnrollLead enrolls a lead in a nurture sequence
func (r *gormNurtureRepository) EnrollLead(enrollment *models.NurtureEnrollment) error {
	return r.db.Create(enrollment).Error
}

// UpdateEnrollment updates a nurture enrollment
func (r *gormNurtureRepository) UpdateEnrollment(enrollment *models.NurtureEnrollment) error {
	return r.db.Save(enrollment).Error
}

// GetEnrollmentActivity returns enrollment activity
func (r *gormNurtureRepository) GetEnrollmentActivity(enrollmentID int) ([]models.NurtureActivity, error) {
	var activities []models.NurtureActivity
	err := r.db.Where("enrollment_id = ?", enrollmentID).Order("created_at desc").Find(&activities).Error
	return activities, err
}

// RecordActivity records nurture activity
func (r *gormNurtureRepository) RecordActivity(activity *models.NurtureActivity) error {
	return r.db.Create(activity).Error
}

// Campaign related methods

// GetCampaigns returns campaigns
func (r *gormNurtureRepository) GetCampaigns(offset int, limit int) ([]models.Campaign, error) {
	var campaigns []models.Campaign
	err := r.db.Offset(offset).Limit(limit).Find(&campaigns).Error
	return campaigns, err
}

// GetCampaignByID returns a campaign by ID
func (r *gormNurtureRepository) GetCampaignByID(id int) (*models.Campaign, error) {
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
func (r *gormNurtureRepository) CreateCampaign(campaign *models.Campaign) error {
	return r.db.Create(campaign).Error
}

// UpdateCampaign updates a campaign
func (r *gormNurtureRepository) UpdateCampaign(campaign *models.Campaign) error {
	return r.db.Save(campaign).Error
}

// DeleteCampaign deletes a campaign
func (r *gormNurtureRepository) DeleteCampaign(id int) error {
	return r.db.Delete(&models.Campaign{}, id).Error
}

// GetCampaignStats returns campaign statistics
func (r *gormNurtureRepository) GetCampaignStats(id int) (map[string]interface{}, error) {
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
func (r *gormNurtureRepository) GetLeadsForCampaign(id int) ([]models.Lead, error) {
	var leads []models.Lead
	err := r.db.Table("lead_campaigns").
		Select("leads.*").
		Joins("JOIN leads ON lead_campaigns.lead_id = leads.id").
		Where("lead_campaigns.campaign_id = ?", id).
		Find(&leads).Error
	return leads, err
}

// AssignLeadsToCampaign assigns leads to a campaign
func (r *gormNurtureRepository) AssignLeadsToCampaign(campaignID int, leadIDs []int) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, leadID := range leadIDs {
		if err := tx.Exec("INSERT INTO lead_campaigns (campaign_id, lead_id) VALUES (?, ?)", campaignID, leadID).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// RemoveLeadsFromCampaign removes leads from a campaign
func (r *gormNurtureRepository) RemoveLeadsFromCampaign(campaignID int, leadIDs []int) error {
	return r.db.Where("campaign_id = ? AND lead_id IN ?", campaignID, leadIDs).Delete("lead_campaigns").Error
}

// GetTemplates returns campaign templates
func (r *gormNurtureRepository) GetTemplates(offset int, limit int) ([]models.CampaignTemplate, error) {
	var templates []models.CampaignTemplate
	err := r.db.Offset(offset).Limit(limit).Find(&templates).Error
	return templates, err
}

// GetTemplateByID returns a campaign template by ID
func (r *gormNurtureRepository) GetTemplateByID(id int) (*models.CampaignTemplate, error) {
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
func (r *gormNurtureRepository) CreateTemplate(template *models.CampaignTemplate) error {
	return r.db.Create(template).Error
}

// UpdateTemplate updates a campaign template
func (r *gormNurtureRepository) UpdateTemplate(template *models.CampaignTemplate) error {
	return r.db.Save(template).Error
}

// DeleteTemplate deletes a campaign template
func (r *gormNurtureRepository) DeleteTemplate(id int) error {
	return r.db.Delete(&models.CampaignTemplate{}, id).Error
}
