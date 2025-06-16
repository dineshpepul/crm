package repositories

import (
	"crm-app/backend/models"
	"time"
)

// type gormAnalyticsRepository struct {
// 	db *gorm.DB
// }

// NewAnalyticsRepository creates a new analytics repository
// func NewAnalyticsRepository(db *gorm.DB) models.AnalyticsRepository {
// 	return &gormAnalyticsRepository{db: db}
// }

// GetLeadAnalytics returns lead analytics with real data from database
func (r *gormAnalyticsRepository) GetLeadAnalytics(startDate time.Time, endDate time.Time, companyId int) (map[string]interface{}, error) {
	var totalLeads int64
	var newLeads int64
	var qualifiedLeads int64

	// Get total leads in date range
	if err := r.db.Model(&models.Lead{}).
		Where("created_at BETWEEN ? AND ? AND company_id = ?", startDate, endDate, companyId).
		Count(&totalLeads).Error; err != nil {
		return nil, err
	}

	// Get new leads (status = 'new')
	if err := r.db.Model(&models.Lead{}).
		Where("created_at BETWEEN ? AND ? AND status = ? AND company_id = ?", startDate, endDate, "new", companyId).
		Count(&newLeads).Error; err != nil {
		return nil, err
	}

	// Get qualified leads (status = 'qualified')
	if err := r.db.Model(&models.Lead{}).
		Where("created_at BETWEEN ? AND ? AND status = ? AND company_id = ?", startDate, endDate, "qualified", companyId).
		Count(&qualifiedLeads).Error; err != nil {
		return nil, err
	}

	// Get leads by source
	var leadsBySource []struct {
		Source string `json:"source"`
		Count  int64  `json:"count"`
	}
	if err := r.db.Model(&models.Lead{}).
		Select("source, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ? AND company_id = ?", startDate, endDate, companyId).
		Group("source").
		Scan(&leadsBySource).Error; err != nil {
		return nil, err
	}

	// Get leads by status
	var leadsByStatus []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	if err := r.db.Model(&models.Lead{}).
		Select("status, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ? AND company_id = ?", startDate, endDate, companyId).
		Group("status").
		Scan(&leadsByStatus).Error; err != nil {
		return nil, err
	}

	// Calculate conversion rate
	var convertedLeads int64
	if err := r.db.Table("deals").
		Joins("JOIN leads ON deals.lead_id = leads.id").
		Where("leads.created_at BETWEEN ? AND ? AND deals.stage = ? AND leads.company_id = ?", startDate, endDate, "won", companyId).
		Count(&convertedLeads).Error; err != nil {
		return nil, err
	}

	conversionRate := float64(0)
	if totalLeads > 0 {
		conversionRate = (float64(convertedLeads) / float64(totalLeads)) * 100
	}

	// Get daily trend data
	var dailyTrend []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	if err := r.db.Model(&models.Lead{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ? AND company_id = ?", startDate, endDate, companyId).
		Group("DATE(created_at)").
		Order("date").
		Scan(&dailyTrend).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_leads":        totalLeads,
		"new_leads":          newLeads,
		"qualified_leads":    qualifiedLeads,
		"leads_by_source":    leadsBySource,
		"leads_by_status":    leadsByStatus,
		"conversion_rate":    conversionRate,
		"average_lead_value": 0, // Would need additional calculation based on deal values
		"daily_trend":        dailyTrend,
	}, nil
}

// GetDealAnalytics returns deal analytics with real data from database
func (r *gormAnalyticsRepository) GetDealAnalytics(startDate time.Time, endDate time.Time, companyId int) (map[string]interface{}, error) {
	var totalDeals int64
	var wonDeals int64
	var lostDeals int64
	var totalRevenue float64
	var avgDealSize float64

	// Get total deals
	if err := r.db.Model(&models.Deal{}).
		Where("created_at BETWEEN ? AND ? AND company_id = ?", startDate, endDate, companyId).
		Count(&totalDeals).Error; err != nil {
		return nil, err
	}

	// Get won deals
	if err := r.db.Model(&models.Deal{}).
		Where("created_at BETWEEN ? AND ? AND stage = ? AND company_id = ?", startDate, endDate, "won", companyId).
		Count(&wonDeals).Error; err != nil {
		return nil, err
	}

	// Get lost deals
	if err := r.db.Model(&models.Deal{}).
		Where("created_at BETWEEN ? AND ? AND stage = ? AND company_id = ?", startDate, endDate, "lost", companyId).
		Count(&lostDeals).Error; err != nil {
		return nil, err
	}

	// Get total revenue from won deals
	if err := r.db.Model(&models.Deal{}).
		Where("created_at BETWEEN ? AND ? AND stage = ? AND company_id = ?", startDate, endDate, "won", companyId).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalRevenue).Error; err != nil {
		return nil, err
	}

	// Calculate average deal size
	if wonDeals > 0 {
		avgDealSize = totalRevenue / float64(wonDeals)
	}

	// Get deals by stage
	var dealsByStage []struct {
		Stage string  `json:"stage"`
		Count int64   `json:"count"`
		Value float64 `json:"value"`
	}
	if err := r.db.Model(&models.Deal{}).
		Select("stage, COUNT(*) as count, COALESCE(SUM(amount), 0) as value").
		Where("created_at BETWEEN ? AND ? AND company_id = ?", startDate, endDate, companyId).
		Group("stage").
		Scan(&dealsByStage).Error; err != nil {
		return nil, err
	}

	// Calculate win rate
	winRate := float64(0)
	closedDeals := wonDeals + lostDeals
	if closedDeals > 0 {
		winRate = (float64(wonDeals) / float64(closedDeals)) * 100
	}

	// Get monthly revenue trend
	var revenueTrend []struct {
		Month   string  `json:"month"`
		Revenue float64 `json:"revenue"`
	}
	if err := r.db.Model(&models.Deal{}).
		Select("DATE_FORMAT(created_at, '%Y-%m') as month, COALESCE(SUM(amount), 0) as revenue").
		Where("created_at BETWEEN ? AND ? AND stage = ? AND company_id = ?", startDate, endDate, "won", companyId).
		Group("DATE_FORMAT(created_at, '%Y-%m')").
		Order("month").
		Scan(&revenueTrend).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_deals":        totalDeals,
		"deals_won":          wonDeals,
		"deals_lost":         lostDeals,
		"deals_by_stage":     dealsByStage,
		"deals_by_value":     dealsByStage, // Same data structure
		"total_revenue":      totalRevenue,
		"average_deal_value": avgDealSize,
		"win_rate":           winRate,
		"deal_velocity":      0, // Would need time-based calculation
		"revenue_trend":      revenueTrend,
	}, nil
}

// GetSalesActivity returns sales activity analytics
func (r *gormAnalyticsRepository) GetSalesActivity(startDate time.Time, endDate time.Time, companyId int) (map[string]interface{}, error) {
	// Note: This would need an activities table to be fully implemented
	// For now, returning calculated data based on leads and deals

	var totalActivities int64

	// Count lead creation as activities
	var leadActivities int64
	if err := r.db.Model(&models.Lead{}).
		Where("created_at BETWEEN ? AND ? AND company_id= ? ", startDate, endDate, companyId).
		Count(&leadActivities).Error; err != nil {
		return nil, err
	}

	// Count deal creation/updates as activities
	var dealActivities int64
	if err := r.db.Model(&models.Deal{}).
		Where("created_at BETWEEN ? AND ? AND company_id= ? ", startDate, endDate, companyId).
		Count(&dealActivities).Error; err != nil {
		return nil, err
	}

	totalActivities = leadActivities + dealActivities

	// Get activities by day
	var activitiesByDay []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}

	// Combine lead and deal activities by day
	query := `
		SELECT DATE(created_at) as date, COUNT(*) as count 
		FROM (
			SELECT created_at FROM leads WHERE created_at BETWEEN ? AND ? AND company_id= ? 
			UNION ALL
			SELECT created_at FROM deals WHERE created_at BETWEEN ? AND ? AND company_id= ? 
		) combined_activities 
		GROUP BY DATE(created_at) 
		ORDER BY date`

	if err := r.db.Raw(query, startDate, endDate, companyId, startDate, endDate, companyId).
		Scan(&activitiesByDay).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_activities":  totalActivities,
		"calls":             0, // Would need activities table
		"emails":            0, // Would need activities table
		"meetings":          0, // Would need activities table
		"notes":             0, // Would need activities table
		"activities_by_day": activitiesByDay,
		"activities_by_type": []map[string]interface{}{
			{"type": "leads", "count": leadActivities},
			{"type": "deals", "count": dealActivities},
		},
	}, nil
}

// GetPerformanceByUser returns performance analytics by user
func (r *gormAnalyticsRepository) GetPerformanceByUser(startDate time.Time, endDate time.Time, companyId int) (map[string]interface{}, error) {
	var userPerformance []struct {
		UserID         uint    `json:"user_id"`
		UserName       string  `json:"user_name"`
		LeadCount      int64   `json:"leads"`
		DealCount      int64   `json:"deals"`
		TotalRevenue   float64 `json:"total_value"`
		ConversionRate float64 `json:"conversion"`
	}

	query := `SELECT l.assigned_to_id AS user_id,COALESCE(l.lead_count,0)AS lead_count,COALESCE(d.deal_count,0)AS deal_count,COALESCE(d.total_revenue,0)AS total_revenue,CASE WHEN COALESCE(l.lead_count,0)>0 THEN(COALESCE(d.deal_count,0)*100.0/l.lead_count)ELSE 0 END AS conversion_rate FROM(SELECT assigned_to_id,COUNT(*)AS lead_count FROM leads WHERE created_at BETWEEN ? AND ? AND company_id= ? AND assigned_to_id IS NOT NULL GROUP BY assigned_to_id)l LEFT JOIN(SELECT assigned_to,COUNT(*)AS deal_count,COALESCE(SUM(amount),0)AS total_revenue FROM deals WHERE created_at BETWEEN ? AND ? AND assigned_to IS NOT NULL AND stage='won' AND company_id= ? GROUP BY assigned_to)d ON l.assigned_to_id=d.assigned_to UNION SELECT d.assigned_to AS user_id,0 AS lead_count,d.deal_count,d.total_revenue,0 AS conversion_rate FROM(SELECT assigned_to,COUNT(*)AS deal_count,COALESCE(SUM(amount),0)AS total_revenue FROM deals WHERE created_at BETWEEN ? AND ? AND assigned_to IS NOT NULL AND stage='won' AND company_id= ? GROUP BY assigned_to)d WHERE d.assigned_to NOT IN(SELECT assigned_to_id FROM leads WHERE created_at BETWEEN ? AND ? AND company_id= ? AND assigned_to_id IS NOT NULL)`

	if err := r.db.Raw(query, startDate, endDate, companyId, startDate, endDate, companyId, startDate, endDate, companyId, startDate, endDate, companyId).
		Scan(&userPerformance).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"users": userPerformance,
	}, nil
}

// GetFunnelAnalytics returns sales funnel analytics
func (r *gormAnalyticsRepository) GetFunnelAnalytics(companyId int) (map[string]interface{}, error) {
	var funnelData []struct {
		Stage string  `json:"stage"`
		Count int64   `json:"count"`
		Value float64 `json:"value"`
	}

	// Get funnel stages with counts and values
	if err := r.db.Model(&models.Deal{}).
		Select("stage, COUNT(*) as count, COALESCE(SUM(amount), 0) as value").
		Group("stage").
		Order("FIELD(stage, 'lead', 'qualified', 'proposal', 'negotiation', 'won', 'lost')").
		Where("company_id = ? ", companyId).
		Scan(&funnelData).Error; err != nil {
		return nil, err
	}

	// Calculate conversion rates between stages
	conversionRates := make([]map[string]interface{}, 0)
	for i := 0; i < len(funnelData)-1; i++ {
		if funnelData[i].Count > 0 {
			rate := (float64(funnelData[i+1].Count) / float64(funnelData[i].Count)) * 100
			conversionRates = append(conversionRates, map[string]interface{}{
				"from_stage": funnelData[i].Stage,
				"to_stage":   funnelData[i+1].Stage,
				"rate":       rate,
			})
		}
	}

	return map[string]interface{}{
		"stages":           funnelData,
		"conversion_rates": conversionRates,
	}, nil
}

// GetTargetAnalytics returns target analytics with dynamic filtering
func (r *gormAnalyticsRepository) GetTargetAnalytics(startDate time.Time, endDate time.Time, userID *uint, companyId int) (map[string]interface{}, error) {
	var targets []struct {
		ID              int       `json:"target_id"`
		Name            string    `json:"name"`
		TargetType      string    `json:"target_type"`
		TargetValue     float64   `json:"target_value"`
		ActualValue     float64   `json:"actual_value"`
		UserId          *int      `json:"user_id"`
		TeamId          *int      `json:"team_id"`
		StartDate       time.Time `json:"start_date"`
		EndDate         time.Time `json:"end_date"`
		Period          string    `json:"period"`
		Status          string    `json:"status"`
		Currency        string    `json:"currency"`
		PercentComplete float64   `json:"percent_complete"`
		TimeProgress    float64   `json:"time_progress"`
		OnTrack         bool      `json:"on_track"`
		DaysRemaining   int       `json:"days_remaining"`
	}

	query := r.db.Model(&models.Target{}).
		Select(`
			id, name, target_type, target_value, actual_value, user_id, team_id,
			start_date, end_date, period, status, currency,
			CASE 
				WHEN target_value > 0 THEN (actual_value / target_value * 100)
				ELSE 0 
			END as percent_complete,
			CASE 
				WHEN DATEDIFF(end_date, start_date) > 0 THEN 
					(DATEDIFF(NOW(), start_date) / DATEDIFF(end_date, start_date) * 100)
				ELSE 0 
			END as time_progress,
			CASE 
				WHEN DATEDIFF(end_date, start_date) > 0 THEN 
					(actual_value / target_value * 100) >= (DATEDIFF(NOW(), start_date) / DATEDIFF(end_date, start_date) * 100)
				ELSE true 
			END as on_track,
			DATEDIFF(end_date, NOW()) as days_remaining
		`).
		Where("start_date <= ? AND end_date >= ? AND company_id = ?", endDate, startDate, companyId)

	// Apply optional filters
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	// Assuming there's a company relationship through users or teams
	query = query.Where("id IN (SELECT id FROM targets)")

	if err := query.Scan(&targets).Error; err != nil {
		return nil, err
	}

	// Calculate actual values for each target based on current data
	for i := range targets {
		actualValue, err := r.calculateActualTargetValue(targets[i].ID, targets[i].TargetType, startDate, endDate, companyId)
		if err == nil {
			targets[i].ActualValue = actualValue
			// Recalculate percent complete with updated actual value
			if targets[i].TargetValue > 0 {
				targets[i].PercentComplete = (actualValue / targets[i].TargetValue) * 100
			}
		}
	}

	// Get targets by type summary
	var targetsByType []struct {
		TargetType  string  `json:"target_type"`
		Count       int64   `json:"count"`
		TotalTarget float64 `json:"total_target"`
		TotalActual float64 `json:"total_actual"`
		AvgProgress float64 `json:"avg_progress"`
	}

	if err := r.db.Model(&models.Target{}).
		Select(`
			target_type,
			COUNT(*) as count,
			SUM(target_value) as total_target,
			SUM(actual_value) as total_actual,
			AVG(CASE WHEN target_value > 0 THEN (actual_value / target_value * 100) ELSE 0 END) as avg_progress
		`).
		Where("start_date <= ? AND end_date >= ? and company_id = ?", endDate, startDate, companyId).
		Group("target_type").
		Scan(&targetsByType).Error; err != nil {
		return nil, err
	}

	// Get targets by period summary
	var targetsByPeriod []struct {
		Period       string  `json:"period"`
		Count        int64   `json:"count"`
		OnTrackCount int64   `json:"on_track_count"`
		AvgProgress  float64 `json:"avg_progress"`
	}

	if err := r.db.Model(&models.Target{}).
		Select(`
			period,
			COUNT(*) as count,
			SUM(CASE WHEN (actual_value / target_value * 100) >= 
				(DATEDIFF(NOW(), start_date) / DATEDIFF(end_date, start_date) * 100) THEN 1 ELSE 0 END) as on_track_count,
			AVG(CASE WHEN target_value > 0 THEN (actual_value / target_value * 100) ELSE 0 END) as avg_progress
		`).
		Where("start_date <= ? AND end_date >= ? and company_id = ?", endDate, startDate, companyId).
		Group("period").
		Scan(&targetsByPeriod).Error; err != nil {
		return nil, err
	}

	// Get monthly progress trend
	var monthlyTrend []struct {
		Month    string  `json:"month"`
		Progress float64 `json:"progress"`
	}

	if err := r.db.Model(&models.Target{}).
		Select(`
			DATE_FORMAT(start_date, '%Y-%m') as month,
			AVG(CASE WHEN target_value > 0 THEN (actual_value / target_value * 100) ELSE 0 END) as progress
		`).
		Where("start_date BETWEEN ? AND ? and company_id = ?", startDate, endDate, companyId).
		Group("DATE_FORMAT(start_date, '%Y-%m')").
		Order("month").
		Scan(&monthlyTrend).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"targets":           targets,
		"targets_by_type":   targetsByType,
		"targets_by_period": targetsByPeriod,
		"monthly_trend":     monthlyTrend,
		"summary": map[string]interface{}{
			"total_targets": len(targets),
			"active_targets": func() int {
				count := 0
				for _, target := range targets {
					if target.Status == "active" {
						count++
					}
				}
				return count
			}(),
		},
	}, nil
}

// Helper method to calculate actual target value based on target type
func (r *gormAnalyticsRepository) calculateActualTargetValue(targetID int, targetType string, startDate, endDate time.Time, companyId int) (float64, error) {
	var actualValue float64

	switch targetType {
	case "revenue":
		// Sum revenue from won deals in the period
		var result struct {
			Value float64
		}
		if err := r.db.Model(&models.Deal{}).
			Select("COALESCE(SUM(amount), 0) as value").
			Where("stage = ? AND created_at BETWEEN ? AND ? AND company_id = ?", "won", startDate, endDate, companyId).
			Scan(&result).Error; err != nil {
			return 0, err
		}
		actualValue = result.Value

	case "leads":
		// Count leads created in the period
		var count int64
		if err := r.db.Model(&models.Lead{}).
			Where("created_at BETWEEN ? AND ? AND company_id = ?", startDate, endDate, companyId).
			Count(&count).Error; err != nil {
			return 0, err
		}
		actualValue = float64(count)

	case "deals":
		// Count deals created in the period
		var count int64
		if err := r.db.Model(&models.Deal{}).
			Where("created_at BETWEEN ? AND ? AND company_id = ?", startDate, endDate, companyId).
			Count(&count).Error; err != nil {
			return 0, err
		}
		actualValue = float64(count)

	case "conversion":
		// Calculate conversion rate
		var leadCount, dealCount int64
		r.db.Model(&models.Lead{}).Where("created_at BETWEEN ? AND ? AND company_id = ?", startDate, endDate, companyId).Count(&leadCount)
		r.db.Model(&models.Deal{}).Where("stage = ? AND created_at BETWEEN ? AND ? AND company_id = ?", "won", startDate, endDate, companyId).Count(&dealCount)

		if leadCount > 0 {
			actualValue = (float64(dealCount) / float64(leadCount)) * 100
		}
	}

	return actualValue, nil
}
