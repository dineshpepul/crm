package handlers

import (
	"net/http"
	"strconv"

	"crm-app/backend/models"

	"github.com/gin-gonic/gin"
)

// CRMNurtureHandler handles requests for lead nurturing
type CRMNurtureHandler struct {
	nurtureRepo models.NurtureRepository
	leadRepo    models.LeadRepository
}

// NewCRMNurtureHandler creates a new nurturing handler
func NewCRMNurtureHandler(repos *models.CRMRepositories) *CRMNurtureHandler {
	return &CRMNurtureHandler{
		nurtureRepo: repos.NurtureRepo,
		leadRepo:    repos.LeadRepo,
	}
}

// GetCampaigns returns all campaigns with pagination
func (h *CRMNurtureHandler) GetCampaigns(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	campaigns, err := h.nurtureRepo.GetCampaigns(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch campaigns"})
		return
	}

	c.JSON(http.StatusOK, campaigns)
}

// GetCampaign returns a campaign by ID
func (h *CRMNurtureHandler) GetCampaign(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	campaign, err := h.nurtureRepo.GetCampaignByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch campaign"})
		return
	}

	if campaign == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}

	c.JSON(http.StatusOK, campaign)
}

// CreateCampaign creates a new campaign
func (h *CRMNurtureHandler) CreateCampaign(c *gin.Context) {
	var campaign models.Campaign
	if err := c.ShouldBindJSON(&campaign); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the created_by field to the current user ID
	// In a real app, would get this from the authenticated user
	userID := 1 // Placeholder
	campaign.CreatedBy = userID

	if err := h.nurtureRepo.CreateCampaign(&campaign); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create campaign"})
		return
	}

	c.JSON(http.StatusCreated, campaign)
}

// UpdateCampaign updates a campaign
func (h *CRMNurtureHandler) UpdateCampaign(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	// Verify that the campaign exists
	existingCampaign, err := h.nurtureRepo.GetCampaignByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch campaign"})
		return
	}
	if existingCampaign == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}

	var campaign models.Campaign
	if err := c.ShouldBindJSON(&campaign); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure ID matches the URL parameter
	campaign.ID = id

	// Preserve the created_by field
	campaign.CreatedBy = existingCampaign.CreatedBy

	if err := h.nurtureRepo.UpdateCampaign(&campaign); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update campaign"})
		return
	}

	c.JSON(http.StatusOK, campaign)
}

// DeleteCampaign deletes a campaign
func (h *CRMNurtureHandler) DeleteCampaign(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	// Verify that the campaign exists
	existingCampaign, err := h.nurtureRepo.GetCampaignByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch campaign"})
		return
	}
	if existingCampaign == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}

	if err := h.nurtureRepo.DeleteCampaign(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete campaign"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campaign deleted successfully"})
}

// GetCampaignStats returns statistics for a campaign
func (h *CRMNurtureHandler) GetCampaignStats(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	stats, err := h.nurtureRepo.GetCampaignStats(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch campaign statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetCampaignLeads returns the leads assigned to a campaign
func (h *CRMNurtureHandler) GetCampaignLeads(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	leads, err := h.nurtureRepo.GetLeadsForCampaign(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch campaign leads"})
		return
	}

	c.JSON(http.StatusOK, leads)
}

// AddLeadsToCampaign adds leads to a campaign
func (h *CRMNurtureHandler) AddLeadsToCampaign(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	// Verify that the campaign exists
	existingCampaign, err := h.nurtureRepo.GetCampaignByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch campaign"})
		return
	}
	if existingCampaign == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}

	var reqBody struct {
		LeadIDs []int `json:"lead_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify that each lead exists
	for _, leadID := range reqBody.LeadIDs {
		lead, err := h.leadRepo.FindByID(leadID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify leads"})
			return
		}
		if lead == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Lead not found",
				"id":    leadID,
			})
			return
		}
	}

	if err := h.nurtureRepo.AssignLeadsToCampaign(id, reqBody.LeadIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add leads to campaign"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Leads added to campaign successfully"})
}

// RemoveLeadsFromCampaign removes leads from a campaign
func (h *CRMNurtureHandler) RemoveLeadsFromCampaign(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	// Verify that the campaign exists
	existingCampaign, err := h.nurtureRepo.GetCampaignByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch campaign"})
		return
	}
	if existingCampaign == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}

	var reqBody struct {
		LeadIDs []int `json:"lead_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.nurtureRepo.RemoveLeadsFromCampaign(id, reqBody.LeadIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove leads from campaign"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Leads removed from campaign successfully"})
}

// GetTemplates returns campaign templates
func (h *CRMNurtureHandler) GetTemplates(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	templates, err := h.nurtureRepo.GetTemplates(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch templates"})
		return
	}

	c.JSON(http.StatusOK, templates)
}

// GetTemplate returns a template by ID
func (h *CRMNurtureHandler) GetTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	template, err := h.nurtureRepo.GetTemplateByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch template"})
		return
	}

	if template == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	c.JSON(http.StatusOK, template)
}

// CreateTemplate creates a new template
func (h *CRMNurtureHandler) CreateTemplate(c *gin.Context) {
	var template models.CampaignTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the created_by field to the current user ID
	userID := 1 // Placeholder
	template.CreatedBy = userID

	if err := h.nurtureRepo.CreateTemplate(&template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create template"})
		return
	}

	c.JSON(http.StatusCreated, template)
}

// UpdateTemplate updates a template
func (h *CRMNurtureHandler) UpdateTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	// Verify that the template exists
	existingTemplate, err := h.nurtureRepo.GetTemplateByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch template"})
		return
	}
	if existingTemplate == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	var template models.CampaignTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure ID matches the URL parameter
	template.ID = id

	// Preserve the created_by field
	template.CreatedBy = existingTemplate.CreatedBy

	if err := h.nurtureRepo.UpdateTemplate(&template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update template"})
		return
	}

	c.JSON(http.StatusOK, template)
}

// DeleteTemplate deletes a template
func (h *CRMNurtureHandler) DeleteTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	// Verify that the template exists
	existingTemplate, err := h.nurtureRepo.GetTemplateByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch template"})
		return
	}
	if existingTemplate == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	if err := h.nurtureRepo.DeleteTemplate(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete template"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template deleted successfully"})
}
