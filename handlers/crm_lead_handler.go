package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"crm-app/backend/models"

	"github.com/gin-gonic/gin"
)

// CRMLeadHandler handles requests for lead management
type CRMLeadHandler struct {
	leadRepo        models.LeadRepository
	fieldConfigRepo models.LeadFieldConfigRepository
}

type CRMScoreHandler struct {
	scoreUpdateRepo models.ScoreRepository
}

// NewCRMLeadHandler creates a new lead handler
func NewCRMLeadHandler(repos *models.CRMRepositories) *CRMLeadHandler {
	return &CRMLeadHandler{
		leadRepo:        repos.LeadRepo,
		fieldConfigRepo: repos.LeadFieldConfigRepo,
	}
}

func NewScoreLeadHandler(repos *models.CRMRepositories) *CRMScoreHandler {
	return &CRMScoreHandler{
		scoreUpdateRepo: repos.LeadScoreType,
	}
}

func (h *CRMScoreHandler) UpdateScore(c *gin.Context) {
	var config []models.ScoreType
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.scoreUpdateRepo.ScoreUpdateRepo(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign lead"})
		return
	}

	c.JSON(http.StatusOK, config)
}

// GetLeads returns all leads
func (h *CRMLeadHandler) GetLeads(c *gin.Context) {
	// Handle query parameters for filtering
	status := c.Query("status")
	assignedToStr := c.Query("assigned_to")
	companyIdStr := c.Query("companyId")
	fmt.Println("companyIdStr", companyIdStr)
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}
	var leads []models.GroupedLead

	// Apply filters if provided
	if status != "" {
		leads, err = h.leadRepo.ListByStatus(status)
	} else if assignedToStr != "" {
		assignedTo, err := strconv.Atoi(assignedToStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assigned_to parameter"})
			return
		}
		leads, err = h.leadRepo.ListByAssignee(assignedTo)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assigned_to parameter"})
			return
		}
	} else {
		leads, err = h.leadRepo.List(companyId)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leads"})
		return
	}

	c.JSON(http.StatusOK, leads)
}

// GetLead returns a lead by ID
func (h *CRMLeadHandler) GetLead(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lead ID"})
		return
	}

	lead, err := h.leadRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch lead"})
		return
	}

	if lead == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lead not found"})
		return
	}

	c.JSON(http.StatusOK, lead)
}

// CreateLead creates a new lead based on dynamic field configuration
func (h *CRMLeadHandler) CreateLead(c *gin.Context) {
	var leadInput models.LeadInput
	if err := c.ShouldBindJSON(&leadInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userId not found in context"})
		return
	}
	userIdStr := userId.(string)
	userIdValue, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userId"})
		return
	}

	lead := models.Lead{

		Status:    "new",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CompanyId: leadInput.CompanyId,
	}
	fmt.Println("hello")
	if err := h.leadRepo.CreateMainLead(&lead); err != nil {
		fmt.Println("check")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	fmt.Println("verify")
	// lastSubmitId, err := h.leadRepo.GetLastSubmitId()
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get last submit id: " + err.Error()})
	// 	return
	// }
	// newSubmitId := lastSubmitId + 1
	newSubmitId := lead.ID
	fmt.Println("leadInput", leadInput)
	var records []models.CrmFieldData
	now := time.Now()
	for _, d := range leadInput.Datas {
		records = append(records, models.CrmFieldData{
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
	fmt.Println("records", records)
	if err := h.leadRepo.Create(records); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create lead: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Records inserted successfully",
		"count":   len(records),
	})
}

// UpdateLead updates a lead
func (h *CRMLeadHandler) UpdateLead(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lead ID"})
		return
	}

	// First, check if the lead exists
	existingLead, err := h.leadRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch lead"})
		return
	}
	if existingLead == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lead not found"})
		return
	}

	// Parse the request body
	var lead models.Lead
	if err := c.ShouldBindJSON(&lead); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure ID matches the URL parameter
	var input uint = uint(id)
	lead.ID = input

	// Update the lead
	if err := h.leadRepo.Update(&lead); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update lead: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, lead)
}

// DeleteLead deletes a lead
func (h *CRMLeadHandler) DeleteLead(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lead ID"})
		return
	}

	// First, check if the lead exists
	existingLead, err := h.leadRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch lead"})
		return
	}
	if existingLead == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lead not found"})
		return
	}

	// Delete the lead
	if err := h.leadRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete lead"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Lead deleted successfully"})
}

// QualifyLead updates a lead's status to qualified
func (h *CRMLeadHandler) QualifyLead(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lead ID"})
		return
	}

	// Get the lead
	lead, err := h.leadRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch lead"})
		return
	}
	if lead == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lead not found"})
		return
	}

	// Update the lead status
	lead.Status = "qualified"

	// Parse score from request if provided
	var reqBody struct {
		Score *int `json:"score"`
	}
	if err := c.ShouldBindJSON(&reqBody); err == nil && reqBody.Score != nil {
		lead.Score = reqBody.Score
	}

	// Update the lead
	if err := h.leadRepo.Update(lead); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to qualify lead"})
		return
	}

	c.JSON(http.StatusOK, lead)
}

// DisqualifyLead updates a lead's status to disqualified
func (h *CRMLeadHandler) DisqualifyLead(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lead ID"})
		return
	}

	// Get the lead
	lead, err := h.leadRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch lead"})
		return
	}
	if lead == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lead not found"})
		return
	}

	// Update the lead status
	lead.Status = "disqualified"

	// Parse reason from request if provided
	var reqBody struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&reqBody); err == nil {
		// Fixed: Initialize CustomFields properly as a slice
		if lead.CustomFields == nil {
			lead.CustomFields = []models.LeadCustomField{}
		}

		// Add disqualification reason as a custom field
		disqualificationReasonField := models.LeadCustomField{
			LeadID:     lead.ID,
			FieldName:  "disqualification_reason",
			FieldValue: reqBody.Reason,
		}

		disqualifiedAtField := models.LeadCustomField{
			LeadID:     lead.ID,
			FieldName:  "disqualified_at",
			FieldValue: time.Now().Format(time.RFC3339),
		}

		lead.CustomFields = append(lead.CustomFields, disqualificationReasonField, disqualifiedAtField)
	}

	// Update the lead
	if err := h.leadRepo.Update(lead); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disqualify lead"})
		return
	}

	c.JSON(http.StatusOK, lead)
}

// AssignLead assigns a lead to a user
func (h *CRMLeadHandler) AssignLead(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lead ID"})
		return
	}

	var reqBody struct {
		AssignedTo int `json:"assigned_to"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the lead
	lead, err := h.leadRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch lead"})
		return
	}
	if lead == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lead not found"})
		return
	}

	// Update the assigned_to field
	assignedIDUint := uint(reqBody.AssignedTo)
	lead.AssignedToID = &assignedIDUint

	// Update the lead
	if err := h.leadRepo.Update(lead); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign lead"})
		return
	}

	c.JSON(http.StatusOK, lead)
}

// func (h *CRMScoreHandler) UpdateScore(c *gin.Context) {
// 	var config []models.ScoreType
// 	if err := c.ShouldBindJSON(&config); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if err := h.scoreUpdateRepo.ScoreUpdateRepo(config); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign lead"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, config)
// }

// BulkImportLeads imports multiple leads from a JSON array
func (h *CRMLeadHandler) BulkImportLeads(c *gin.Context) {
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

	lead := models.Lead{
		Status:    "new",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CompanyId: 1,
	}

	if err := h.leadRepo.CreateMainLead(&lead); err != nil {
		fmt.Println("check")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	newSubmitId := lead.ID

	// lastSubmitId, err := h.leadRepo.GetLastSubmitId()
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get last submit id: " + err.Error()})
	// 	return
	// }
	// newSubmitId := lastSubmitId + 1

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

	c.JSON(http.StatusCreated, gin.H{
		"message": "Bulk records inserted successfully",
		"count":   len(allRecords),
	})
}

// ExportLeads exports all leads or filtered leads
func (h *CRMLeadHandler) ExportLeads(c *gin.Context) {
	// Handle query parameters for filtering
	status := c.Query("status")
	assignedToStr := c.Query("assigned_to")
	companyIdStr := c.Query("companyId")
	companyId, err1 := strconv.Atoi(companyIdStr)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}

	var leads []models.GroupedLead
	var err error

	// Apply filters if provided
	if status != "" {
		leads, err = h.leadRepo.ListByStatus(status)
	} else if assignedToStr != "" {
		assignedTo, err := strconv.Atoi(assignedToStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assigned_to parameter"})
			return
		}
		leads, err = h.leadRepo.ListByAssignee(assignedTo)
	} else {
		leads, err = h.leadRepo.List(companyId)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leads for export"})
		return
	}

	// Get all field configurations for complete data export
	fieldConfigs, err := h.fieldConfigRepo.GetAllFieldConfigs(companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch field configurations"})
		return
	}

	// Create export data structure with field definitions and data
	exportData := gin.H{
		"generated_at":      time.Now().Format(time.RFC3339),
		"total_leads":       len(leads),
		"field_definitions": fieldConfigs,
		"leads":             leads,
	}

	c.JSON(http.StatusOK, exportData)
}

// GetAllFormSections returns all form sections
func (h *CRMLeadHandler) GetAllFormSections(c *gin.Context) {
	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}
	sections, err := h.fieldConfigRepo.GetAllFormSections(companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch form sections"})
		return
	}

	c.JSON(http.StatusOK, sections)
}

// GetVisibleFormSections returns visible form sections
func (h *CRMLeadHandler) GetVisibleFormSections(c *gin.Context) {
	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}
	sections, err := h.fieldConfigRepo.GetVisibleFormSections(companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch visible form sections"})
		return
	}

	c.JSON(http.StatusOK, sections)
}

// CreateFormSection creates a new form section
func (h *CRMLeadHandler) CreateFormSection(c *gin.Context) {
	var section models.LeadFormSection
	if err := c.ShouldBindJSON(&section); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.fieldConfigRepo.CreateFormSection(&section); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create form section"})
		return
	}

	c.JSON(http.StatusCreated, section)
}

// UpdateFormSection updates a form section
func (h *CRMLeadHandler) UpdateFormSection(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
		return
	}

	var section models.LeadFormSection
	if err := c.ShouldBindJSON(&section); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure ID matches the URL parameter
	section.ID = uint(id)

	if err := h.fieldConfigRepo.UpdateFormSection(&section); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update form section"})
		return
	}

	c.JSON(http.StatusOK, section)
}

// DeleteFormSection deletes a form section
func (h *CRMLeadHandler) DeleteFormSection(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
		return
	}

	if err := h.fieldConfigRepo.DeleteFormSection(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete form section"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Form section deleted successfully"})
}

// ReorderFormSections updates the order of form sections
func (h *CRMLeadHandler) ReorderFormSections(c *gin.Context) {
	var reqBody struct {
		SectionIDs []int `json:"section_ids"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.fieldConfigRepo.ReorderFormSections(reqBody.SectionIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reorder form sections"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Form sections reordered successfully"})
}
