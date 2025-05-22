
package models

import (
	"time"
)

// AnalyticsPeriod represents a time period for analytics
type AnalyticsPeriod struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// LeadAnalytics represents analytics data for leads
type LeadAnalytics struct {
	TotalLeads       int                      `json:"total_leads"`
	NewLeads         int                      `json:"new_leads"`
	QualifiedLeads   int                      `json:"qualified_leads"`
	LeadsBySource    map[string]int           `json:"leads_by_source"`
	LeadsByStatus    map[string]int           `json:"leads_by_status"`
	ConversionRate   float64                  `json:"conversion_rate"`
	LeadGrowthByDay  map[string]int           `json:"lead_growth_by_day"`
	SourceEfficiency map[string]float64       `json:"source_efficiency"`
	LeadQuality      map[string][]interface{} `json:"lead_quality"`
}

// DealAnalytics represents analytics data for deals
type DealAnalytics struct {
	TotalDeals          int                `json:"total_deals"`
	DealsWon            int                `json:"deals_won"`
	DealsLost           int                `json:"deals_lost"`
	TotalRevenue        float64            `json:"total_revenue"`
	AverageDealSize     float64            `json:"average_deal_size"`
	AverageSalesCycle   int                `json:"average_sales_cycle"` // in days
	DealsByStage        map[string]int     `json:"deals_by_stage"`
	RevenueTrend        map[string]float64 `json:"revenue_trend"`
	WinRate             float64            `json:"win_rate"`
	DealVelocity        float64            `json:"deal_velocity"`
	RevenueForecast     float64            `json:"revenue_forecast"`
	DealAgeDistribution map[string]int     `json:"deal_age_distribution"`
}

// SalesActivityAnalytics represents analytics data for sales activities
type SalesActivityAnalytics struct {
	TotalActivities       int                `json:"total_activities"`
	ActivitiesByType      map[string]int     `json:"activities_by_type"`
	ActivitiesByUser      map[string]int     `json:"activities_by_user"`
	ActivitiesByStage     map[string]int     `json:"activities_by_stage"`
	ActivitiesTrend       map[string]int     `json:"activities_trend"`
	AvgActivitiesPerDeal  float64            `json:"avg_activities_per_deal"`
	AvgActivitiesPerStage map[string]float64 `json:"avg_activities_per_stage"`
	ActivityEfficiency    map[string]float64 `json:"activity_efficiency"`
}

// UserPerformanceAnalytics represents analytics data for user performance
type UserPerformanceAnalytics struct {
	LeadsGenerated     map[string]int     `json:"leads_generated"`
	LeadsConverted     map[string]int     `json:"leads_converted"`
	DealsWon           map[string]int     `json:"deals_won"`
	RevenueGenerated   map[string]float64 `json:"revenue_generated"`
	AvgDealSize        map[string]float64 `json:"avg_deal_size"`
	AvgDealCycle       map[string]int     `json:"avg_deal_cycle"`
	ActivityCompletion map[string]float64 `json:"activity_completion"`
	PerformanceRank    map[string]int     `json:"performance_rank"`
	ConversionRate     map[string]float64 `json:"conversion_rate"`
}

// FunnelAnalytics represents analytics data for the sales funnel
type FunnelAnalytics struct {
	StageNames            []string             `json:"stage_names"`
	StageCount            map[string]int       `json:"stage_count"`
	ConversionRates       map[string]float64   `json:"conversion_rates"`
	AvgTimeInStage        map[string]float64   `json:"avg_time_in_stage"` // in days
	DropOffPoints         map[string]int       `json:"drop_off_points"`
	StageVelocity         map[string]float64   `json:"stage_velocity"`
	StageEfficiency       map[string]float64   `json:"stage_efficiency"`
	StageTrends           map[string][]int     `json:"stage_trends"`
	FunnelHealth          string               `json:"funnel_health"` // "healthy", "at_risk", "critical"
	FunnelRecommendations []string             `json:"funnel_recommendations"`
	BottleneckStages      []map[string]string  `json:"bottleneck_stages"`
}

// Note: The AnalyticsRepository interface is now only defined in repository_interfaces.go
// and the implementation will be in repositories/analytics_repository.go
