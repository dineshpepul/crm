package handlers

import (
	"net/http"
	"strconv"
	"time"

	"crm-app/backend/models"
	"crm-app/backend/services"

	"github.com/gin-gonic/gin"
)

// LeadHandler handles lead-related requests
type LeadHandler struct {
	leadService *services.LeadService
	leadRepo    models.LeadRepository
}

// NewLeadHandler creates a new LeadHandler
func NewLeadHandler(leadService *services.LeadService) *LeadHandler {
	return &LeadHandler{
		leadService: leadService,
	}
}

// GetLeads returns all leads with optional filtering
func (h *LeadHandler) GetLeads(c *gin.Context) {
	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}
	leads, err := h.leadService.GetLeads(companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leads: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, leads)
}

// GetLead returns a specific lead by ID
func (h *LeadHandler) GetLead(c *gin.Context) {
	id, err := h.getIDParam(c)
	if err != nil {
		return
	}

	lead, err := h.leadService.GetLeadByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch lead: " + err.Error()})
		return
	}

	if lead == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lead not found"})
		return
	}

	c.JSON(http.StatusOK, lead)
}

// CreateLead creates a new lead
func (h *LeadHandler) CreateLead(c *gin.Context) {
	var bulkInput []models.LeadInput
	if err := c.ShouldBindJSON(&bulkInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(bulkInput) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No leads provided for import"})
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userId not found in context"})
		return
	}
	userIdValue := userId.(int)

	var allRecords []models.CrmFieldData
	now := time.Now()

	lastSubmitId, err := h.leadRepo.GetLastSubmitId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get last submit id: " + err.Error()})
		return
	}
	newSubmitId := lastSubmitId + 1

	for _, leadInput := range bulkInput {
		for _, d := range leadInput.Datas {
			allRecords = append(allRecords, models.CrmFieldData{
				CompanyId:  leadInput.CompanyId,
				CrmStageId: d.StageId,
				CrmFieldId: d.FieldId,
				FieldValue: d.FieldValue,
				CreatedBy:  userIdValue,
				SubmitId:   newSubmitId,
				CreatedAt:  now,
				UpdatedAt:  now,
			})
		}
	}
	if err := h.leadRepo.Create(allRecords); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create leads: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, allRecords)
}

// UpdateLead updates an existing lead
func (h *LeadHandler) UpdateLead(c *gin.Context) {
	id, err := h.getIDParam(c)
	if err != nil {
		return
	}

	// Check if the lead exists
	existingLead, err := h.leadService.GetLeadByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch lead: " + err.Error()})
		return
	}
	if existingLead == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lead not found"})
		return
	}

	var leadUpdate struct {
		Name         string                 `json:"name"`
		Email        string                 `json:"email"`
		Phone        string                 `json:"phone"`
		Company      string                 `json:"company"`
		Source       string                 `json:"source"`
		Status       string                 `json:"status"`
		Notes        string                 `json:"notes"`
		Tags         []string               `json:"tags"`
		CustomFields map[string]interface{} `json:"custom_fields"`
	}

	if err := c.ShouldBindJSON(&leadUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Update fields
	existingLead.Name = leadUpdate.Name
	existingLead.Email = leadUpdate.Email
	existingLead.Phone = leadUpdate.Phone
	existingLead.Company = leadUpdate.Company
	existingLead.Source = leadUpdate.Source
	existingLead.Status = leadUpdate.Status
	existingLead.Notes = leadUpdate.Notes
	existingLead.Tags = leadUpdate.Tags

	// Update custom fields
	existingLead.CustomFields = nil
	for field, value := range leadUpdate.CustomFields {
		existingLead.CustomFields = append(existingLead.CustomFields, models.LeadCustomField{
			LeadID:     existingLead.ID,
			FieldName:  field,
			FieldValue: toString(value),
		})
	}

	if err := h.leadService.UpdateLead(existingLead); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update lead: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, existingLead)
}

// DeleteLead deletes a lead
func (h *LeadHandler) DeleteLead(c *gin.Context) {
	id, err := h.getIDParam(c)
	if err != nil {
		return
	}

	// Check if the lead exists
	existingLead, err := h.leadService.GetLeadByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch lead: " + err.Error()})
		return
	}
	if existingLead == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lead not found"})
		return
	}

	if err := h.leadService.DeleteLead(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete lead: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Lead deleted successfully"})
}

// QualifyLead marks a lead as qualified
func (h *LeadHandler) QualifyLead(c *gin.Context) {
	id, err := h.getIDParam(c)
	if err != nil {
		return
	}

	// Parse score from request if provided
	var reqBody struct {
		Score *int `json:"score"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		// Ignore binding errors, just proceed with nil score
		reqBody.Score = nil
	}

	if err := h.leadService.QualifyLead(id, reqBody.Score); err != nil {
		if err == services.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Lead not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to qualify lead: " + err.Error()})
		}
		return
	}

	lead, _ := h.leadService.GetLeadByID(id)
	c.JSON(http.StatusOK, lead)
}

// DisqualifyLead marks a lead as disqualified
func (h *LeadHandler) DisqualifyLead(c *gin.Context) {
	id, err := h.getIDParam(c)
	if err != nil {
		return
	}

	if err := h.leadService.DisqualifyLead(id); err != nil {
		if err == services.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Lead not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disqualify lead: " + err.Error()})
		}
		return
	}

	lead, _ := h.leadService.GetLeadByID(id)
	c.JSON(http.StatusOK, lead)
}

// AssignLead assigns a lead to a user
func (h *LeadHandler) AssignLead(c *gin.Context) {
	id, err := h.getIDParam(c)
	if err != nil {
		return
	}

	var reqBody struct {
		AssignedToID int `json:"assigned_to_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if err := h.leadService.AssignLead(id, reqBody.AssignedToID); err != nil {
		if err == services.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Lead not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign lead: " + err.Error()})
		}
		return
	}

	lead, _ := h.leadService.GetLeadByID(id)
	c.JSON(http.StatusOK, lead)
}

// GetAllFieldConfigs returns all field configurations
func (h *LeadHandler) GetAllFieldConfigs(c *gin.Context) {
	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}
	configs, err := h.leadService.GetAllFieldConfigs(companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch field configurations: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, configs)
}

// GetVisibleFieldConfigs returns visible field configurations
func (h *LeadHandler) GetVisibleFieldConfigs(c *gin.Context) {
	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}
	configs, err := h.leadService.GetVisibleFieldConfigs(companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch visible field configurations: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, configs)
}

// GetRequiredFieldConfigs returns required field configurations
func (h *LeadHandler) GetRequiredFieldConfigs(c *gin.Context) {
	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}
	configs, err := h.leadService.GetRequiredFieldConfigs(companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch required field configurations: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, configs)
}

// GetFieldConfigsBySection returns field configurations by section
func (h *LeadHandler) GetFieldConfigsBySection(c *gin.Context) {
	section := c.Param("section")
	if section == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Section parameter is required"})
		return
	}

	configs, err := h.leadService.GetFieldConfigsBySection(section)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch field configurations: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, configs)
}

// CreateFieldConfig creates a new field configuration
func (h *LeadHandler) CreateFieldConfig(c *gin.Context) {
	var config models.LeadFieldConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if err := h.leadService.CreateFieldConfig(&config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create field configuration: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, config)
}

// UpdateFieldConfig updates a field configuration
func (h *LeadHandler) UpdateFieldConfig(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var config models.LeadFieldConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	config.ID = uint(id)

	if err := h.leadService.UpdateFieldConfig(&config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update field configuration: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

// DeleteFieldConfig deletes a field configuration
func (h *LeadHandler) DeleteFieldConfig(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.leadService.DeleteFieldConfig(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete field configuration: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Field configuration deleted successfully"})
}

// ReorderFormFields updates the order of form fields
func (h *LeadHandler) ReorderFormFields(c *gin.Context) {
	var request struct {
		FieldIDs []uint `json:"field_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if err := h.leadService.ReorderFormFields(request.FieldIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reorder fields: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Fields reordered successfully"})
}

// BulkImportLeads imports multiple leads
func (h *LeadHandler) BulkImportLeads(c *gin.Context) {
	var request struct {
		Leads []models.Lead `json:"leads" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	result, err := h.leadService.BulkImportLeads(request.Leads)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import leads: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ExportLeads exports leads
func (h *LeadHandler) ExportLeads(c *gin.Context) {
	// Get query parameters as filters
	filters := make(map[string]string)
	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			filters[key] = values[0]
		}
	}

	leads, err := h.leadService.ExportLeads(filters, companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export leads: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, leads)
}

// GetAllFormSections returns all form sections
func (h *LeadHandler) GetAllFormSections(c *gin.Context) {
	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}
	sections, err := h.leadService.GetAllFormSections(companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch form sections: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, sections)
}

// GetVisibleFormSections returns visible form sections
func (h *LeadHandler) GetVisibleFormSections(c *gin.Context) {
	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}
	sections, err := h.leadService.GetVisibleFormSections(companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch visible form sections: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, sections)
}

// CreateFormSection creates a new form section
func (h *LeadHandler) CreateFormSection(c *gin.Context) {
	var section models.LeadFormSection
	if err := c.ShouldBindJSON(&section); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if err := h.leadService.CreateFormSection(&section); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create form section: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, section)
}

// UpdateFormSection updates a form section
func (h *LeadHandler) UpdateFormSection(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var section models.LeadFormSection
	if err := c.ShouldBindJSON(&section); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	section.ID = uint(id)

	if err := h.leadService.UpdateFormSection(&section); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update form section: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, section)
}

// DeleteFormSection deletes a form section
func (h *LeadHandler) DeleteFormSection(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.leadService.DeleteFormSection(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete form section: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Form section deleted successfully"})
}

// ReorderFormSections updates the order of form sections
func (h *LeadHandler) ReorderFormSections(c *gin.Context) {
	var request struct {
		SectionIDs []int `json:"section_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if err := h.leadService.ReorderFormSections(request.SectionIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reorder sections: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Form sections reordered successfully"})
}

// Helper function to get the ID parameter from the URL
func (h *LeadHandler) getIDParam(c *gin.Context) (int, error) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return 0, err
	}
	return id, nil
}

// Helper function to convert any value to string
func toString(value interface{}) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case int, int64, float64, bool:
		return strconv.FormatInt(int64(v.(float64)), 10)
	default:
		return ""
	}
}
