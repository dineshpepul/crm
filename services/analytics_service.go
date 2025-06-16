package services

import (
	"crm-app/backend/models"
	"fmt"
	"time"
)

// AnalyticsService handles business logic for analytics
type AnalyticsService struct {
	analyticsRepo models.AnalyticsRepository
	leadRepo      models.LeadRepository
	dealRepo      models.DealRepository
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(repos *models.CRMRepositories) *AnalyticsService {
	return &AnalyticsService{
		analyticsRepo: repos.AnalyticsRepo,
		leadRepo:      repos.LeadRepo,
		dealRepo:      repos.DealRepo,
	}
}

// AnalyticsFilters represents filters for analytics queries
type AnalyticsFilters struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CompanyId int       `json:"company_id"`
	UserID    *uint     `json:"user_id,omitempty"`
	Source    string    `json:"source,omitempty"`
	Status    string    `json:"status,omitempty"`
	Stage     string    `json:"stage,omitempty"`
	Campaign  string    `json:"campaign,omitempty"`
}

// GetLeadAnalytics retrieves lead analytics with business logic
func (s *AnalyticsService) GetLeadAnalytics(filters AnalyticsFilters, companyId int) (map[string]interface{}, error) {
	// Validate date range
	if filters.EndDate.Before(filters.StartDate) {
		return nil, fmt.Errorf("end date cannot be before start date")
	}

	// Get raw analytics data
	analytics, err := s.analyticsRepo.GetLeadAnalytics(filters.StartDate, filters.EndDate, companyId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch lead analytics: %w", err)
	}

	// Add calculated metrics
	analytics["metrics"] = s.calculateLeadMetrics(analytics)

	// Add growth indicators
	previousPeriod := s.calculatePreviousPeriod(filters.StartDate, filters.EndDate, companyId)
	prevAnalytics, _ := s.analyticsRepo.GetLeadAnalytics(previousPeriod.StartDate, previousPeriod.EndDate, companyId)
	analytics["growth"] = s.calculateGrowthMetrics(analytics, prevAnalytics)

	return analytics, nil
}

// GetDealAnalytics retrieves deal analytics with business logic
func (s *AnalyticsService) GetDealAnalytics(filters AnalyticsFilters, companyId int) (map[string]interface{}, error) {
	// Validate date range
	if filters.EndDate.Before(filters.StartDate) {
		return nil, fmt.Errorf("end date cannot be before start date")
	}

	analytics, err := s.analyticsRepo.GetDealAnalytics(filters.StartDate, filters.EndDate, companyId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch deal analytics: %w", err)
	}

	// Add calculated metrics
	analytics["metrics"] = s.calculateDealMetrics(analytics)

	// Add forecasting
	analytics["forecast"] = s.calculateRevenueForecast(analytics)

	return analytics, nil
}

// GetSalesActivityAnalytics retrieves sales activity analytics
func (s *AnalyticsService) GetSalesActivityAnalytics(filters AnalyticsFilters, companyId int) (map[string]interface{}, error) {
	analytics, err := s.analyticsRepo.GetSalesActivity(filters.StartDate, filters.EndDate, companyId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sales activity analytics: %w", err)
	}

	// Add productivity metrics
	analytics["productivity"] = s.calculateProductivityMetrics(analytics)

	return analytics, nil
}

// GetPerformanceAnalytics retrieves user performance analytics
func (s *AnalyticsService) GetPerformanceAnalytics(filters AnalyticsFilters, companyId int) (map[string]interface{}, error) {
	analytics, err := s.analyticsRepo.GetPerformanceByUser(filters.StartDate, filters.EndDate, companyId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch performance analytics: %w", err)
	}

	// Add performance rankings
	analytics["rankings"] = s.calculatePerformanceRankings(analytics)

	return analytics, nil
}

// GetFunnelAnalytics retrieves sales funnel analytics
func (s *AnalyticsService) GetFunnelAnalytics(companyId int) (map[string]interface{}, error) {
	analytics, err := s.analyticsRepo.GetFunnelAnalytics(companyId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch funnel analytics: %w", err)
	}

	// Add funnel health score
	analytics["health_score"] = s.calculateFunnelHealthScore(analytics)

	return analytics, nil
}

// GetDashboardAnalytics provides comprehensive dashboard data
func (s *AnalyticsService) GetDashboardAnalytics(filters AnalyticsFilters, companyId int) (map[string]interface{}, error) {
	leadAnalytics, err := s.GetLeadAnalytics(filters, companyId)
	if err != nil {
		return nil, err
	}

	dealAnalytics, err := s.GetDealAnalytics(filters, companyId)
	if err != nil {
		return nil, err
	}

	dashboard := map[string]interface{}{
		"period": map[string]interface{}{
			"start_date": filters.StartDate.Format("2006-01-02"),
			"end_date":   filters.EndDate.Format("2006-01-02"),
		},
		"summary": map[string]interface{}{
			"total_leads":        leadAnalytics["total_leads"],
			"qualified_leads":    leadAnalytics["qualified_leads"],
			"conversion_rate":    leadAnalytics["conversion_rate"],
			"total_deals":        dealAnalytics["total_deals"],
			"deals_won":          dealAnalytics["deals_won"],
			"total_revenue":      dealAnalytics["total_revenue"],
			"average_deal_value": dealAnalytics["average_deal_value"],
			"win_rate":           dealAnalytics["win_rate"],
		},
		"charts": map[string]interface{}{
			"leads_by_source": leadAnalytics["leads_by_source"],
			"deals_by_stage":  dealAnalytics["deals_by_stage"],
			"revenue_trend":   dealAnalytics["revenue_trend"],
			"daily_trend":     leadAnalytics["daily_trend"],
		},
		"insights": s.generateInsights(leadAnalytics, dealAnalytics),
	}

	return dashboard, nil
}

// Helper methods for calculations
func (s *AnalyticsService) calculateLeadMetrics(analytics map[string]interface{}) map[string]interface{} {
	totalLeads, _ := analytics["total_leads"].(int64)
	qualifiedLeads, _ := analytics["qualified_leads"].(int64)

	qualificationRate := float64(0)
	if totalLeads > 0 {
		qualificationRate = (float64(qualifiedLeads) / float64(totalLeads)) * 100
	}

	return map[string]interface{}{
		"qualification_rate": qualificationRate,
		"lead_velocity":      s.calculateLeadVelocity(analytics),
		"source_efficiency":  s.calculateSourceEfficiency(analytics),
	}
}

func (s *AnalyticsService) calculateDealMetrics(analytics map[string]interface{}) map[string]interface{} {
	totalRevenue, _ := analytics["total_revenue"].(float64)
	dealsWon, _ := analytics["deals_won"].(int64)

	return map[string]interface{}{
		"revenue_per_deal": func() float64 {
			if dealsWon > 0 {
				return totalRevenue / float64(dealsWon)
			}
			return 0
		}(),
		"deal_velocity": s.calculateDealVelocity(analytics),
	}
}

func (s *AnalyticsService) calculateProductivityMetrics(analytics map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"activities_per_day": s.calculateActivitiesPerDay(analytics),
		"efficiency_score":   s.calculateEfficiencyScore(analytics),
	}
}

func (s *AnalyticsService) calculatePerformanceRankings(analytics map[string]interface{}) map[string]interface{} {
	// Implementation for performance rankings
	return map[string]interface{}{
		"top_performers":            []map[string]interface{}{},
		"improvement_opportunities": []map[string]interface{}{},
	}
}

func (s *AnalyticsService) calculateFunnelHealthScore(analytics map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"score":  75, // Example score
		"status": "healthy",
		"recommendations": []string{
			"Focus on lead qualification",
			"Improve conversion at proposal stage",
		},
	}
}

func (s *AnalyticsService) calculatePreviousPeriod(startDate, endDate time.Time, companyId int) AnalyticsFilters {
	duration := endDate.Sub(startDate)
	return AnalyticsFilters{
		StartDate: startDate.Add(-duration),
		EndDate:   startDate,
	}
}

func (s *AnalyticsService) calculateGrowthMetrics(current, previous map[string]interface{}) map[string]interface{} {
	currentLeads, _ := current["total_leads"].(int64)
	previousLeads, _ := previous["total_leads"].(int64)

	growthRate := float64(0)
	if previousLeads > 0 {
		growthRate = ((float64(currentLeads) - float64(previousLeads)) / float64(previousLeads)) * 100
	}

	return map[string]interface{}{
		"leads_growth_rate": growthRate,
		"trend": func() string {
			if growthRate > 0 {
				return "up"
			} else if growthRate < 0 {
				return "down"
			}
			return "stable"
		}(),
	}
}

func (s *AnalyticsService) GetTargetAnalytics(filters AnalyticsFilters) (map[string]interface{}, error) {
	// Validate date range
	if filters.EndDate.Before(filters.StartDate) {
		return nil, fmt.Errorf("end date cannot be before start date")
	}

	// Get target analytics data
	analytics, err := s.analyticsRepo.GetTargetAnalytics(filters.StartDate, filters.EndDate, filters.UserID, filters.CompanyId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch target analytics: %w", err)
	}

	// Add calculated metrics and insights
	analytics["insights"] = s.generateTargetInsights(analytics)
	analytics["performance_summary"] = s.calculateTargetPerformanceSummary(analytics)

	// Add previous period comparison
	previousPeriod := s.calculatePreviousPeriod(filters.StartDate, filters.EndDate, filters.CompanyId)
	prevAnalytics, _ := s.analyticsRepo.GetTargetAnalytics(previousPeriod.StartDate, previousPeriod.EndDate, filters.UserID, filters.CompanyId)
	analytics["comparison"] = s.calculateTargetComparison(analytics, prevAnalytics)

	return analytics, nil
}

func (s *AnalyticsService) calculateTargetComparison(current, previous map[string]interface{}) map[string]interface{} {
	comparison := map[string]interface{}{
		"targets_growth":     0.0,
		"performance_change": "stable",
	}

	currentTargets := 0
	previousTargets := 0

	if targets, ok := current["targets"].([]map[string]interface{}); ok {
		currentTargets = len(targets)
	}

	if targets, ok := previous["targets"].([]map[string]interface{}); ok {
		previousTargets = len(targets)
	}

	if previousTargets > 0 {
		growth := ((float64(currentTargets) - float64(previousTargets)) / float64(previousTargets)) * 100
		comparison["targets_growth"] = growth

		if growth > 5 {
			comparison["performance_change"] = "improving"
		} else if growth < -5 {
			comparison["performance_change"] = "declining"
		}
	}

	return comparison
}

func (s *AnalyticsService) generateTargetInsights(analytics map[string]interface{}) []string {
	insights := []string{}

	if targets, ok := analytics["targets"].([]map[string]interface{}); ok {
		onTrackCount := 0
		totalTargets := len(targets)

		for _, target := range targets {
			if onTrack, ok := target["on_track"].(bool); ok && onTrack {
				onTrackCount++
			}
		}

		if totalTargets > 0 {
			onTrackPercentage := (float64(onTrackCount) / float64(totalTargets)) * 100
			if onTrackPercentage < 50 {
				insights = append(insights, "Less than 50% of targets are on track. Consider reviewing target settings or increasing team focus.")
			} else if onTrackPercentage >= 80 {
				insights = append(insights, "Excellent performance! 80% or more of targets are on track.")
			}
		}
	}

	return insights
}

func (s *AnalyticsService) calculateTargetPerformanceSummary(analytics map[string]interface{}) map[string]interface{} {
	summary := map[string]interface{}{
		"total_targets":           0,
		"on_track":                0,
		"behind":                  0,
		"achieved":                0,
		"overall_completion_rate": 0.0,
	}

	if targets, ok := analytics["targets"].([]map[string]interface{}); ok {
		totalTargets := len(targets)
		summary["total_targets"] = totalTargets

		onTrackCount := 0
		achievedCount := 0
		totalCompletion := 0.0

		for _, target := range targets {
			if completion, ok := target["percent_complete"].(float64); ok {
				totalCompletion += completion
				if completion >= 100 {
					achievedCount++
				}
			}

			if onTrack, ok := target["on_track"].(bool); ok && onTrack {
				onTrackCount++
			}
		}

		summary["on_track"] = onTrackCount
		summary["behind"] = totalTargets - onTrackCount
		summary["achieved"] = achievedCount

		if totalTargets > 0 {
			summary["overall_completion_rate"] = totalCompletion / float64(totalTargets)
		}
	}

	return summary
}

func (s *AnalyticsService) calculateRevenueForecast(analytics map[string]interface{}) map[string]interface{} {
	// Simple forecasting based on current trends
	totalRevenue, _ := analytics["total_revenue"].(float64)

	return map[string]interface{}{
		"next_month":   totalRevenue * 1.1, // 10% growth assumption
		"next_quarter": totalRevenue * 3.2,
		"confidence":   0.75,
	}
}

// Additional helper methods
func (s *AnalyticsService) calculateLeadVelocity(analytics map[string]interface{}) float64 {
	// Implementation for lead velocity calculation
	return 0.0
}

func (s *AnalyticsService) calculateSourceEfficiency(analytics map[string]interface{}) map[string]interface{} {
	// Implementation for source efficiency calculation
	return map[string]interface{}{}
}

func (s *AnalyticsService) calculateDealVelocity(analytics map[string]interface{}) float64 {
	// Implementation for deal velocity calculation
	return 0.0
}

func (s *AnalyticsService) calculateActivitiesPerDay(analytics map[string]interface{}) float64 {
	// Implementation for activities per day calculation
	return 0.0
}

func (s *AnalyticsService) calculateEfficiencyScore(analytics map[string]interface{}) float64 {
	// Implementation for efficiency score calculation
	return 0.0
}

func (s *AnalyticsService) generateInsights(leadAnalytics, dealAnalytics map[string]interface{}) []string {
	insights := []string{}

	conversionRate, _ := leadAnalytics["conversion_rate"].(float64)
	if conversionRate < 10 {
		insights = append(insights, "Lead conversion rate is below average. Consider improving lead qualification.")
	}

	winRate, _ := dealAnalytics["win_rate"].(float64)
	if winRate < 20 {
		insights = append(insights, "Deal win rate needs improvement. Focus on better prospect targeting.")
	}

	return insights
}
