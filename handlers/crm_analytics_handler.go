package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"crm-app/backend/models"
	"crm-app/backend/services"

	"github.com/gin-gonic/gin"
)

// CRMAnalyticsHandler handles requests for CRM analytics
type CRMAnalyticsHandler struct {
	analyticsService *services.AnalyticsService
}

// NewCRMAnalyticsHandler creates a new analytics handler
func NewCRMAnalyticsHandler(repos *models.CRMRepositories) *CRMAnalyticsHandler {
	return &CRMAnalyticsHandler{
		analyticsService: services.NewAnalyticsService(repos),
	}
}

// GetLeadAnalytics returns lead analytics for a date range
func (h *CRMAnalyticsHandler) GetLeadAnalytics(c *gin.Context) {
	filters, err := h.parseAnalyticsFilters(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}
	analytics, err := h.analyticsService.GetLeadAnalytics(filters, companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch lead analytics"})
		return
	}

	response := map[string]interface{}{
		"period": map[string]interface{}{
			"start_date": filters.StartDate.Format("2006-01-02"),
			"end_date":   filters.EndDate.Format("2006-01-02"),
		},
		"filters": map[string]interface{}{
			"company_id": filters.CompanyId,
			"user_id":    filters.UserID,
			"source":     filters.Source,
			"status":     filters.Status,
		},
		"data": analytics,
	}

	c.JSON(http.StatusOK, response)
}

// GetDealAnalytics returns deal analytics for a date range
func (h *CRMAnalyticsHandler) GetDealAnalytics(c *gin.Context) {
	filters, err := h.parseAnalyticsFilters(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}

	analytics, err := h.analyticsService.GetDealAnalytics(filters, companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch deal analytics"})
		return
	}

	response := map[string]interface{}{
		"period": map[string]interface{}{
			"start_date": filters.StartDate.Format("2006-01-02"),
			"end_date":   filters.EndDate.Format("2006-01-02"),
		},
		"filters": map[string]interface{}{
			"company_id": filters.CompanyId,
			"user_id":    filters.UserID,
			"stage":      filters.Stage,
		},
		"data": analytics,
	}

	c.JSON(http.StatusOK, response)
}

// GetSalesActivityAnalytics returns sales activity analytics
func (h *CRMAnalyticsHandler) GetSalesActivityAnalytics(c *gin.Context) {
	filters, err := h.parseAnalyticsFilters(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}

	analytics, err := h.analyticsService.GetSalesActivityAnalytics(filters, companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sales activity analytics"})
		return
	}

	response := map[string]interface{}{
		"period": map[string]interface{}{
			"start_date": filters.StartDate.Format("2006-01-02"),
			"end_date":   filters.EndDate.Format("2006-01-02"),
		},
		"data": analytics,
	}

	c.JSON(http.StatusOK, response)
}

// GetPerformanceAnalytics returns performance analytics by user
func (h *CRMAnalyticsHandler) GetPerformanceAnalytics(c *gin.Context) {
	filters, err := h.parseAnalyticsFilters(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}

	analytics, err := h.analyticsService.GetPerformanceAnalytics(filters, companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch performance analytics"})
		return
	}

	response := map[string]interface{}{
		"period": map[string]interface{}{
			"start_date": filters.StartDate.Format("2006-01-02"),
			"end_date":   filters.EndDate.Format("2006-01-02"),
		},
		"data": analytics,
	}

	c.JSON(http.StatusOK, response)
}

// GetFunnelAnalytics returns sales funnel analytics
func (h *CRMAnalyticsHandler) GetFunnelAnalytics(c *gin.Context) {
	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}
	analytics, err := h.analyticsService.GetFunnelAnalytics(companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch funnel analytics"})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetTargetAnalytics returns target analytics
// func (h *CRMAnalyticsHandler) GetTargetAnalytics(c *gin.Context) {
// 	period := c.DefaultQuery("period", "month")

// 	// For now, return placeholder data as targets would need a separate implementation
// 	// You would implement this based on your targets table structure
// 	targetData := map[string]interface{}{
// 		"target":           100000,
// 		"achieved":         75000,
// 		"achievement_rate": 75,
// 		"remaining":        25000,
// 		"forecast":         95000,
// 		"forecast_percent": 95,
// 	}

// 	prevPeriodData := map[string]interface{}{
// 		"target":           90000,
// 		"achieved":         85000,
// 		"achievement_rate": 94,
// 	}

// 	c.JSON(http.StatusOK, map[string]interface{}{
// 		"period":          period,
// 		"data":            targetData,
// 		"previous_period": prevPeriodData,
// 	})
// }

// GetDashboardAnalytics returns a combined analytics dashboard
func (h *CRMAnalyticsHandler) GetDashboardAnalytics(c *gin.Context) {
	filters, err := h.parseAnalyticsFilters(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}

	dashboard, err := h.analyticsService.GetDashboardAnalytics(filters, companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch dashboard analytics"})
		return
	}

	c.JSON(http.StatusOK, dashboard)
}

// GetConversionAnalytics returns lead-to-deal conversion analytics
func (h *CRMAnalyticsHandler) GetConversionAnalytics(c *gin.Context) {
	filters, err := h.parseAnalyticsFilters(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}

	leadAnalytics, err := h.analyticsService.GetLeadAnalytics(filters, companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch conversion analytics"})
		return
	}

	// Extract conversion rate from lead analytics
	overallConversionRate := leadAnalytics["conversion_rate"]

	conversionAnalytics := map[string]interface{}{
		"period": map[string]interface{}{
			"start_date": filters.StartDate.Format("2006-01-02"),
			"end_date":   filters.EndDate.Format("2006-01-02"),
		},
		"overall_conversion_rate": overallConversionRate,
		"leads_by_source":         leadAnalytics["leads_by_source"],
		"leads_by_status":         leadAnalytics["leads_by_status"],
		"conversion_funnel":       leadAnalytics["metrics"],
	}

	c.JSON(http.StatusOK, conversionAnalytics)
}

// Helper method to parse analytics filters from query parameters
func (h *CRMAnalyticsHandler) parseAnalyticsFilters(c *gin.Context) (services.AnalyticsFilters, error) {
	var filters services.AnalyticsFilters

	// Parse date range
	startDateStr := c.DefaultQuery("start_date", "")
	endDateStr := c.DefaultQuery("end_date", "")

	if startDateStr == "" {
		// Default to last 30 days
		filters.StartDate = time.Now().AddDate(0, 0, -30)
	} else {
		var err error
		filters.StartDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return filters, fmt.Errorf("invalid start_date format. Use YYYY-MM-DD")
		}
	}

	if endDateStr == "" {
		// Default to today
		filters.EndDate = time.Now()
	} else {
		var err error
		filters.EndDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return filters, fmt.Errorf("invalid end_date format. Use YYYY-MM-DD")
		}
	}

	// Adjust endDate to include the entire day
	filters.EndDate = filters.EndDate.Add(24*time.Hour - time.Second)

	// Parse optional filters
	// if companyIDStr := c.Query("company_id"); companyIDStr != "" {
	// 	if companyId, err := strconv.ParseUint(companyIDStr, 10, 32); err == nil {
	// 		filters.CompanyId = int(companyId)
	// 	}
	// }

	if companyIDStr := c.Query("companyId"); companyIDStr != "" {
		if companyId, err := strconv.Atoi(companyIDStr); err == nil {
			filters.CompanyId = companyId
		}
	}

	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			id := uint(userID)
			filters.UserID = &id
		}
	}

	filters.Source = c.Query("source")
	filters.Status = c.Query("status")
	filters.Stage = c.Query("stage")
	filters.Campaign = c.Query("campaign")

	return filters, nil
}

func (h *CRMAnalyticsHandler) GetTargetAnalytics(c *gin.Context) {
	filters, err := h.parseAnalyticsFilters(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	analytics, err := h.analyticsService.GetTargetAnalytics(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch target analytics"})
		return
	}

	response := map[string]interface{}{
		"period": map[string]interface{}{
			"start_date": filters.StartDate.Format("2006-01-02"),
			"end_date":   filters.EndDate.Format("2006-01-02"),
		},
		"filters": map[string]interface{}{
			"company_id":  filters.CompanyId,
			"user_id":     filters.UserID,
			"target_type": c.Query("target_type"),
			"period":      c.Query("period"),
		},
		"data": analytics,
	}

	c.JSON(http.StatusOK, response)
}
