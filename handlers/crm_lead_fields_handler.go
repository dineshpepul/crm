package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"crm-app/backend/models"

	"github.com/gin-gonic/gin"
)

// CRMLeadFieldsHandler handles lead field configurations
type CRMLeadFieldsHandler struct {
	fieldConfigRepo models.LeadFieldConfigRepository
}

// NewCRMLeadFieldsHandler creates a new lead fields handler
func NewCRMLeadFieldsHandler(repos *models.CRMRepositories) *CRMLeadFieldsHandler {
	return &CRMLeadFieldsHandler{
		fieldConfigRepo: repos.LeadFieldConfigRepo,
	}
}

// GetAllFieldConfigs returns all field configs
func (h *CRMLeadFieldsHandler) GetAllFieldConfigs(c *gin.Context) {
	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}
	configs, err := h.fieldConfigRepo.GetAllFieldConfigs(companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch field configurations"})
		return
	}

	c.JSON(http.StatusOK, configs)
}

// GetVisibleFieldConfigs returns visible field configs
func (h *CRMLeadFieldsHandler) GetVisibleFieldConfigs(c *gin.Context) {
	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}
	configs, err := h.fieldConfigRepo.GetVisibleFieldConfigs(companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch visible field configurations"})
		return
	}

	c.JSON(http.StatusOK, configs)
}

// GetRequiredFieldConfigs returns required field configs
func (h *CRMLeadFieldsHandler) GetRequiredFieldConfigs(c *gin.Context) {
	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}
	configs, err := h.fieldConfigRepo.GetRequiredFieldConfigs(companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch required field configurations"})
		return
	}

	c.JSON(http.StatusOK, configs)
}

// GetFieldConfigsBySection returns field configs for a section
func (h *CRMLeadFieldsHandler) GetFieldConfigsBySection(c *gin.Context) {
	section := c.Param("section")

	configs, err := h.fieldConfigRepo.GetFieldConfigsBySection(section)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch field configurations for section"})
		return
	}

	c.JSON(http.StatusOK, configs)
}

// GetFieldConfig returns a field config by ID
func (h *CRMLeadFieldsHandler) GetFieldConfig(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid field config ID"})
		return
	}

	config, err := h.fieldConfigRepo.GetFieldConfig(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch field configuration"})
		return
	}

	if config == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Field configuration not found"})
		return
	}

	c.JSON(http.StatusOK, config)
}

// CreateFieldConfig creates a new field config
func (h *CRMLeadFieldsHandler) CreateFieldConfig(c *gin.Context) {
	var config models.LeadFieldConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.fieldConfigRepo.CreateFieldConfig(&config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create field configuration"})
		return
	}

	c.JSON(http.StatusCreated, config)
}

// UpdateFieldConfig updates a field config
func (h *CRMLeadFieldsHandler) UpdateFieldConfig(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid field config ID"})
		return
	}

	var config models.LeadFieldConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.ID = uint(id)

	if err := h.fieldConfigRepo.UpdateFieldConfig(&config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update field configuration"})
		return
	}

	c.JSON(http.StatusOK, config)
}

// DeleteFieldConfig deletes a field config
func (h *CRMLeadFieldsHandler) DeleteFieldConfig(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid field config ID"})
		return
	}

	if err := h.fieldConfigRepo.DeleteFieldConfig(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete field configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Field configuration deleted successfully"})
}

// ReorderFormFields reorders form fields
func (h *CRMLeadFieldsHandler) ReorderFormFields(c *gin.Context) {
	var request struct {
		FieldIDs []int `json:"field_ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert int slice to uint slice
	var uintFieldIDs []uint
	for _, id := range request.FieldIDs {
		uintFieldIDs = append(uintFieldIDs, uint(id))
	}

	if err := h.fieldConfigRepo.ReorderFormFields(uintFieldIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reorder form fields"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Form fields reordered successfully"})
}

// GetAllFormSections returns all form sections
func (h *CRMLeadFieldsHandler) GetAllFormSections(c *gin.Context) {
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
func (h *CRMLeadFieldsHandler) GetVisibleFormSections(c *gin.Context) {
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
func (h *CRMLeadFieldsHandler) CreateFormSection(c *gin.Context) {
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
func (h *CRMLeadFieldsHandler) UpdateFormSection(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form section ID"})
		return
	}

	var section models.LeadFormSection
	if err := c.ShouldBindJSON(&section); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	section.ID = uint(id)

	if err := h.fieldConfigRepo.UpdateFormSection(&section); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update form section"})
		return
	}

	c.JSON(http.StatusOK, section)
}

// DeleteFormSection deletes a form section
func (h *CRMLeadFieldsHandler) DeleteFormSection(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form section ID"})
		return
	}

	if err := h.fieldConfigRepo.DeleteFormSection(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete form section"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Form section deleted successfully"})
}

// ReorderFormSections reorders form sections
func (h *CRMLeadFieldsHandler) ReorderFormSections(c *gin.Context) {
	var request struct {
		SectionIDs []int `json:"section_ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.fieldConfigRepo.ReorderFormSections(request.SectionIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reorder form sections"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Form sections reordered successfully"})
}

// GetFormStructure returns the complete form structure
func (h *CRMLeadFieldsHandler) GetFormStructure(c *gin.Context) {
	companyIdStr := c.Query("companyId")
	fmt.Println("companyIdStr", companyIdStr)
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}
	structure, err := h.fieldConfigRepo.GetFormStructure(companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch form structure"})
		return
	}

	c.JSON(http.StatusOK, structure)
}
