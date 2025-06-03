package repositories

import (
	"crm-app/backend/models"
	"time"
)

// Implementation of DashboardRepository that uses GORM
// type gormDashboardRepository struct {
// 	db *gorm.DB
// }

// GetDashboardSummary retrieves summary statistics for the dashboard
func (r *gormDashboardRepository) GetDashboardSummary(companyId int) (map[string]interface{}, error) {
	summary := make(map[string]interface{})

	// Get total leads
	var totalLeads int64
	if err := r.db.Model(&models.Lead{}).Where("company_id = ?", companyId).Count(&totalLeads).Error; err != nil {
		return nil, err
	}
	summary["total_leads"] = totalLeads

	// Get new leads today
	var newLeadsToday int64
	today := time.Now().Format("2006-01-02")
	if err := r.db.Model(&models.Lead{}).Where("DATE(created_at) = ? AND company_id = ?", today, companyId).Count(&newLeadsToday).Error; err != nil {
		return nil, err
	}
	summary["new_leads_today"] = newLeadsToday

	// Get qualified leads
	var qualifiedLeads int64
	if err := r.db.Model(&models.Lead{}).Where("status = ? AND company_id = ?", "qualified", companyId).Count(&qualifiedLeads).Error; err != nil {
		return nil, err
	}
	summary["qualified_leads"] = qualifiedLeads

	// Get total deals, won deals, lost deals
	var totalDeals, dealsWon, dealsLost int64
	if err := r.db.Model(&models.Deal{}).Where("company_id = ?", companyId).Count(&totalDeals).Error; err != nil {
		return nil, err
	}
	summary["total_deals"] = totalDeals

	if err := r.db.Model(&models.Deal{}).Where("stage = ? AND company_id = ?", "won", companyId).Count(&dealsWon).Error; err != nil {
		return nil, err
	}
	summary["deals_won"] = dealsWon

	if err := r.db.Model(&models.Deal{}).Where("stage = ? AND company_id = ?", "lost", companyId).Count(&dealsLost).Error; err != nil {
		return nil, err
	}
	summary["deals_lost"] = dealsLost

	// Calculate conversion rate
	if totalLeads > 0 {
		conversionRate := float64(dealsWon) / float64(totalLeads) * 100
		summary["conversion_rate"] = conversionRate
	} else {
		summary["conversion_rate"] = 0.0
	}

	// Get revenue metrics
	var totalRevenue float64
	if err := r.db.Model(&models.Deal{}).Where("stage = ? AND company_id = ?", "won", companyId).Select("COALESCE(SUM(amount), 0)").Row().Scan(&totalRevenue); err != nil {
		return nil, err
	}
	summary["total_revenue"] = totalRevenue

	var forecastedRevenue float64
	if err := r.db.Model(&models.Deal{}).Where("stage NOT IN (?, ?) AND company_id = ?", "won", "lost", companyId).
		Select("COALESCE(SUM(amount * probability / 100), 0)").Row().Scan(&forecastedRevenue); err != nil {
		return nil, err
	}
	summary["forecasted_revenue"] = forecastedRevenue

	// Get average deal size
	var avgDealSize float64
	if err := r.db.Model(&models.Deal{}).Where("amount > 0").Select("COALESCE(AVG(amount), 0)").Row().Scan(&avgDealSize); err != nil {
		return nil, err
	}
	summary["average_deal_size"] = avgDealSize

	// Get average sales cycle
	// DATEDIFF in MySQL calculates days between dates
	var avgSalesCycle float64
	err := r.db.Raw(`
		SELECT COALESCE(AVG(DATEDIFF(deals.created_at, leads.created_at)), 0) 
		FROM deals 
		JOIN leads ON deals.lead_id = leads.id 
		WHERE deals.stage = 'won' AND leads.company_id=?
	`, companyId).Row().Scan(&avgSalesCycle)
	if err != nil {
		return nil, err
	}
	summary["average_sales_cycle"] = int(avgSalesCycle)

	summary["last_updated"] = time.Now()

	return summary, nil
}

// GetLeadsBySource retrieves lead counts grouped by source
func (r *gormDashboardRepository) GetLeadsBySource(companyId int) ([]map[string]interface{}, error) {
	var results []struct {
		Source string
		Count  int
	}

	if err := r.db.Model(&models.Lead{}).
		Select("source, COUNT(*) as count").
		Where("company_id = ?", companyId).
		Group("source").
		Order("count DESC").
		Find(&results).Error; err != nil {
		return nil, err
	}

	var output []map[string]interface{}
	for _, result := range results {
		output = append(output, map[string]interface{}{
			"source": result.Source,
			"count":  result.Count,
		})
	}

	return output, nil
}

// GetLeadsByStatus retrieves lead counts grouped by status
func (r *gormDashboardRepository) GetLeadsByStatus(companyId int) ([]map[string]interface{}, error) {
	var results []struct {
		Status string
		Count  int
	}

	if err := r.db.Model(&models.Lead{}).
		Select("status, COUNT(*) as count").
		Where("company_id = ?", companyId).
		Group("status").
		Order("count DESC").
		Find(&results).Error; err != nil {
		return nil, err
	}

	var output []map[string]interface{}
	for _, result := range results {
		output = append(output, map[string]interface{}{
			"status": result.Status,
			"count":  result.Count,
		})
	}

	return output, nil
}

// GetRevenueByMonth retrieves monthly revenue for a given year
func (r *gormDashboardRepository) GetRevenueByMonth(year int, companyId int) ([]map[string]interface{}, error) {
	var results []struct {
		Month   int
		Revenue float64
	}

	if err := r.db.Model(&models.Deal{}).
		Select("MONTH(created_at) as month, SUM(amount) as revenue").
		Where("YEAR(created_at) = ? AND stage = ? AND company_id= ?", year, "won", companyId).
		Group("MONTH(created_at)").
		Order("month").
		Find(&results).Error; err != nil {
		return nil, err
	}

	var output []map[string]interface{}
	for _, result := range results {
		// Convert month number to month name
		monthName := time.Month(result.Month).String()
		output = append(output, map[string]interface{}{
			"month":   monthName,
			"revenue": result.Revenue,
		})
	}

	return output, nil
}

// GetSalesForecast retrieves sales forecast for the coming months
func (r *gormDashboardRepository) GetSalesForecast(months int, companyId int) ([]map[string]interface{}, error) {
	var results []struct {
		Year             int
		Month            int
		ForecastedAmount float64
	}

	// This is a simplified forecast based on probability-weighted deals
	if err := r.db.Model(&models.Deal{}).
		Select("YEAR(expected_close_date) as year, MONTH(expected_close_date) as month, SUM(amount * probability / 100) as forecasted_amount").
		Where("company_id = ? AND stage NOT IN (?, ?) AND expected_close_date IS NOT NULL AND expected_close_date <= DATE_ADD(CURDATE(), INTERVAL ? MONTH)",
			companyId, "won", "lost", months).
		Group("YEAR(expected_close_date), MONTH(expected_close_date)").
		Order("year, month").
		Find(&results).Error; err != nil {
		return nil, err
	}

	var output []map[string]interface{}
	for _, result := range results {
		// Format as YYYY-MM
		date := time.Date(result.Year, time.Month(result.Month), 1, 0, 0, 0, 0, time.UTC)
		key := date.Format("2006-01")
		output = append(output, map[string]interface{}{
			"period": key,
			"amount": result.ForecastedAmount,
		})
	}

	return output, nil
}
