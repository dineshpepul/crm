package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"crm-app/backend/models"

	"github.com/gin-gonic/gin"
)

// CRMDealHandler handles requests for deal management
type CRMDealHandler struct {
	dealRepo models.DealRepository
	leadRepo models.LeadRepository
}

// NewCRMDealHandler creates a new deal handler
func NewCRMDealHandler(repos *models.CRMRepositories) *CRMDealHandler {
	return &CRMDealHandler{
		dealRepo: repos.DealRepo,
		leadRepo: repos.LeadRepo,
	}
}

// GetDeals returns all deals with optional filtering
func (h *CRMDealHandler) GetDeals(c *gin.Context) {
	// Parse query parameters for filtering
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	stage := c.Query("stage")
	assignedToStr := c.Query("assigned_to")
	minAmountStr := c.Query("amount")
	maxAmountStr := c.Query("amount")
	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}

	filters := make(map[string]interface{})

	if stage != "" {
		filters["stage"] = stage
	}

	if assignedToStr != "" {
		assignedTo, err := strconv.Atoi(assignedToStr)
		if err == nil {
			filters["assigned_to"] = assignedTo
		}
	}

	if minAmountStr != "" {
		minAmount, err := strconv.ParseFloat(minAmountStr, 64)
		if err == nil {
			filters["amount"] = minAmount
		}
	}

	if maxAmountStr != "" {
		maxAmount, err := strconv.ParseFloat(maxAmountStr, 64)
		if err == nil {
			filters["amount"] = maxAmount
		}
	}

	deals, err := h.dealRepo.List(offset, limit, filters, companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch deals"})
		return
	}

	c.JSON(http.StatusOK, deals)
}

// GetDeal returns a deal by ID
func (h *CRMDealHandler) GetDeal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deal ID"})
		return
	}

	deal, err := h.dealRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch deal"})
		return
	}

	if deal == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deal not found"})
		return
	}

	c.JSON(http.StatusOK, deal)
}

// CreateDeal creates a new deal
func (h *CRMDealHandler) CreateDeal(c *gin.Context) {
	var deal models.Deal
	if err := c.ShouldBindJSON(&deal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify that the lead exists
	lead, err := h.leadRepo.FindByID(deal.LeadID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify lead"})
		return
	}
	if lead == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lead does not exist"})
		return
	}

	if err := h.dealRepo.Create(&deal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create deal"})
		return
	}

	c.JSON(http.StatusCreated, deal)
}

// UpdateDeal updates a deal
func (h *CRMDealHandler) UpdateDeal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deal ID"})
		return
	}

	// Verify that the deal exists
	existingDeal, err := h.dealRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch deal"})
		return
	}
	if existingDeal == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deal not found"})
		return
	}

	var deal models.Deal
	if err := c.ShouldBindJSON(&deal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure ID matches the URL parameter
	deal.ID = id
	fmt.Println("")
	// Verify that the lead exists
	lead, err := h.leadRepo.FindByID(deal.LeadID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify lead"})
		return
	}
	if lead == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lead does not exist"})
		return
	}

	if err := h.dealRepo.Update(&deal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update deal"})
		return
	}

	c.JSON(http.StatusOK, deal)
}

// DeleteDeal deletes a deal
func (h *CRMDealHandler) DeleteDeal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deal ID"})
		return
	}

	// Verify that the deal exists
	existingDeal, err := h.dealRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch deal"})
		return
	}
	if existingDeal == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deal not found"})
		return
	}

	if err := h.dealRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete deal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deal deleted successfully"})
}

// GetDealsByLead returns deals for a specific lead
func (h *CRMDealHandler) GetDealsByLead(c *gin.Context) {
	leadIDStr := c.Param("lead_id")
	leadID, err := strconv.Atoi(leadIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lead ID"})
		return
	}

	deals, err := h.dealRepo.FindByLead(leadID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch deals"})
		return
	}

	c.JSON(http.StatusOK, deals)
}

// GetDealPipeline returns deal pipeline analytics
func (h *CRMDealHandler) GetDealPipeline(c *gin.Context) {
	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}
	pipeline, err := h.dealRepo.GetDealPipeline(companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch deal pipeline"})
		return
	}

	c.JSON(http.StatusOK, pipeline)
}

// UpdateDealStage updates a deal's stage
func (h *CRMDealHandler) UpdateDealStage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deal ID"})
		return
	}

	var reqBody struct {
		Stage string `json:"stage" binding:"required"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the existing deal
	deal, err := h.dealRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch deal"})
		return
	}
	if deal == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deal not found"})
		return
	}

	// Update the stage
	deal.Stage = reqBody.Stage

	// Auto-update probability based on stage
	switch reqBody.Stage {
	case "prospecting":
		deal.Probability = 10
	case "qualification":
		deal.Probability = 25
	case "needs_analysis":
		deal.Probability = 40
	case "proposal":
		deal.Probability = 60
	case "negotiation":
		deal.Probability = 80
	case "won":
		deal.Probability = 100
	case "lost":
		deal.Probability = 0
	}

	if err := h.dealRepo.Update(deal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update deal stage"})
		return
	}

	c.JSON(http.StatusOK, deal)
}
