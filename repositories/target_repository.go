package repositories

import (
	"crm-app/backend/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

// gormTargetRepository implements TargetRepository with GORM
// type gormTargetRepository struct {
// 	db *gorm.DB
// }

// NewTargetRepository creates a new target repository
// func NewTargetRepository(db *gorm.DB) models.TargetRepository {
// 	return &gormTargetRepository{db: db}
// }

// GetTargets returns targets with filters
func (r *gormTargetRepository) GetTargets(filters map[string]interface{}, companyId int) ([]models.Target, error) {
	var targets []models.Target
	query := r.db

	// Apply filters if provided
	if filters != nil {
		for key, value := range filters {
			query = query.Where(key+" = ?", value)
		}
	}

	if err := query.Where("company_id=?", companyId).Find(&targets).Error; err != nil {
		return nil, err
	}
	return targets, nil
}

// GetTargetByID gets a target by ID
func (r *gormTargetRepository) GetTargetByID(id int) (*models.Target, error) {
	var target models.Target
	result := r.db.First(&target, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &target, nil
}

// CreateTarget creates a new target
func (r *gormTargetRepository) CreateTarget(target *models.Target) error {
	return r.db.Create(target).Error
}

// UpdateTarget updates an existing target
func (r *gormTargetRepository) UpdateTarget(target *models.Target) error {
	return r.db.Omit("CreatedAt").Save(target).Error
}

// DeleteTarget deletes a target
func (r *gormTargetRepository) DeleteTarget(id int) error {
	return r.db.Delete(&models.Target{}, id).Error
}

// GetTargetProgress gets the progress toward a target
func (r *gormTargetRepository) GetTargetProgress(id int, companyId int) (map[string]interface{}, error) {
	var target models.Target
	if err := r.db.Where("company_id=?", companyId).First(&target, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Calculate progress based on target type
	var actualValue float64

	switch target.TargetType {
	case "revenue":
		// Sum the revenue from closed-won deals in the target period
		var result struct {
			Value float64
		}
		if err := r.db.Model(&models.Deal{}).
			Select("COALESCE(SUM(amount), 0) as value").
			Where("company_id = ? AND stage = ? AND created_at BETWEEN ? AND ?", companyId, "closed_won", target.StartDate, target.EndDate).
			Scan(&result).Error; err != nil {
			return nil, err
		}
		actualValue = result.Value

	case "leads":
		// Count leads created in the target period
		var count int64
		if err := r.db.Model(&models.Lead{}).
			Where("company_id = ? AND created_at BETWEEN ? AND ?", companyId, target.StartDate, target.EndDate).
			Count(&count).Error; err != nil {
			return nil, err
		}
		actualValue = float64(count)

	case "deals":
		// Count deals created in the target period
		var count int64
		if err := r.db.Model(&models.Deal{}).
			Where("company_id = ? AND created_at BETWEEN ? AND ?", companyId, target.StartDate, target.EndDate).
			Count(&count).Error; err != nil {
			return nil, err
		}
		actualValue = float64(count)
	}

	// Update the actual value in the database
	if err := r.db.Model(&target).Where("company_id = ?", companyId).Update("actual_value", actualValue).Error; err != nil {
		return nil, err
	}

	// Calculate percentage of target achieved
	percentComplete := (actualValue / target.TargetValue) * 100

	// Calculate time progress
	now := time.Now()
	var timeProgress float64

	if now.Before(target.StartDate) {
		timeProgress = 0
	} else if now.After(target.EndDate) {
		timeProgress = 100
	} else {
		totalDuration := target.EndDate.Sub(target.StartDate)
		elapsedDuration := now.Sub(target.StartDate)
		timeProgress = (float64(elapsedDuration) / float64(totalDuration)) * 100
	}

	// Calculate if target is on track
	var onTrack bool
	if timeProgress <= 0 {
		onTrack = true // Target period hasn't started yet
	} else {
		onTrack = percentComplete >= timeProgress
	}

	return map[string]interface{}{
		"target_id":        target.ID,
		"target_type":      target.TargetType,
		"target_value":     target.TargetValue,
		"actual_value":     actualValue,
		"percent_complete": percentComplete,
		"time_progress":    timeProgress,
		"start_date":       target.StartDate,
		"end_date":         target.EndDate,
		"on_track":         onTrack,
		"period":           target.Period,
	}, nil
}

// GetAllTargetProgress gets progress for all active targets
func (r *gormTargetRepository) GetAllTargetProgress(companyId int) ([]map[string]interface{}, error) {
	var targets []models.Target
	if err := r.db.Where("status = ? AND company_id=?", "active", companyId).Find(&targets).Error; err != nil {
		return nil, err
	}

	progress := make([]map[string]interface{}, 0, len(targets))

	for _, target := range targets {
		targetProgress, err := r.GetTargetProgress(target.ID, companyId)
		if err != nil {
			return nil, err
		}

		progress = append(progress, targetProgress)
	}

	return progress, nil
}
