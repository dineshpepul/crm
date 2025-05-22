
package models

import (
	"time"
)

// DashboardSummary represents summary statistics for the CRM dashboard
type DashboardSummary struct {
	TotalLeads          int       `json:"total_leads"`
	NewLeadsToday       int       `json:"new_leads_today"`
	QualifiedLeads      int       `json:"qualified_leads"`
	ConversionRate      float64   `json:"conversion_rate"`
	TotalDeals          int       `json:"total_deals"`
	DealsWon            int       `json:"deals_won"`
	DealsLost           int       `json:"deals_lost"`
	TotalRevenue        float64   `json:"total_revenue"`
	ForecastedRevenue   float64   `json:"forecasted_revenue"`
	AverageDealSize     float64   `json:"average_deal_size"`
	AverageSalesCycle   int       `json:"average_sales_cycle"` // in days
	LastUpdated         time.Time `json:"last_updated"`
}

// LeadsBySource represents lead count grouped by source
type LeadsBySource struct {
	Source string `json:"source"`
	Count  int    `json:"count"`
}

// LeadsByStatus represents lead count grouped by status
type LeadsByStatus struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

// Note: The DashboardRepository interface is now only defined in repository_interfaces.go
// and the implementation will be in repositories/dashboard_repository.go
