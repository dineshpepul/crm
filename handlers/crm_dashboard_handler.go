package handlers

import (
	"net/http"
	"strconv"
	"time"

	"crm-app/backend/models"

	"github.com/gin-gonic/gin"
)

// CRMDashboardHandler handles requests for the CRM dashboard
type CRMDashboardHandler struct {
	dashboardRepo models.DashboardRepository
	leadRepo      models.LeadRepository
	dealRepo      models.DealRepository
	targetRepo    models.TargetRepository
}

// NewCRMDashboardHandler creates a new dashboard handler
func NewCRMDashboardHandler(repos *models.CRMRepositories) *CRMDashboardHandler {
	return &CRMDashboardHandler{
		dashboardRepo: repos.DashboardRepo,
		leadRepo:      repos.LeadRepo,
		dealRepo:      repos.DealRepo,
		targetRepo:    repos.TargetRepo,
	}
}

// GetDashboardSummary returns summary statistics for the dashboard
func (h *CRMDashboardHandler) GetDashboardSummary(c *gin.Context) {
	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}
	summary, err := h.dashboardRepo.GetDashboardSummary(companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch dashboard summary"})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetLeadsBySource returns lead count by source
func (h *CRMDashboardHandler) GetLeadsBySource(c *gin.Context) {
	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}
	leads, err := h.dashboardRepo.GetLeadsBySource(companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leads by source"})
		return
	}

	c.JSON(http.StatusOK, leads)
}

// GetLeadsByStatus returns lead count by status
func (h *CRMDashboardHandler) GetLeadsByStatus(c *gin.Context) {
	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}
	leads, err := h.dashboardRepo.GetLeadsByStatus(companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leads by status"})
		return
	}

	c.JSON(http.StatusOK, leads)
}

// GetRevenueByMonth returns revenue by month for the current year
func (h *CRMDashboardHandler) GetRevenueByMonth(c *gin.Context) {
	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}
	year := time.Now().Year()
	revenue, err := h.dashboardRepo.GetRevenueByMonth(year, companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch revenue by month"})
		return
	}

	c.JSON(http.StatusOK, revenue)
}

// GetSalesForecast returns sales forecast for the next 6 months
func (h *CRMDashboardHandler) GetSalesForecast(c *gin.Context) {
	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}
	forecast, err := h.dashboardRepo.GetSalesForecast(6, companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sales forecast"})
		return
	}

	c.JSON(http.StatusOK, forecast)
}

// GetTopDeals returns the top deals by amount
func (h *CRMDashboardHandler) GetTopDeals(c *gin.Context) {
	limit := 5
	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}
	// Get deals sorted by amount, limiting to 5 results
	deals, err := h.dealRepo.List(0, limit, map[string]interface{}{}, companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch top deals"})
		return
	}

	c.JSON(http.StatusOK, deals)
}

// GetRecentLeads returns the most recent leads
func (h *CRMDashboardHandler) GetRecentLeads(c *gin.Context) {
	limit := 5
	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}
	// Get leads sorted by created_at desc, limiting to 5 results
	leads, err := h.leadRepo.List(companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recent leads"})
		return
	}

	// Take only the first few leads
	if len(leads) > limit {
		leads = leads[:limit]
	}

	c.JSON(http.StatusOK, leads)
}

// GetTargetProgress returns progress towards active sales targets
func (h *CRMDashboardHandler) GetTargetProgress(c *gin.Context) {
	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}
	progress, err := h.targetRepo.GetAllTargetProgress(companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch target progress"})
		return
	}

	c.JSON(http.StatusOK, progress)
}
