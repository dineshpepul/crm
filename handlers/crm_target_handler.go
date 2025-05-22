package handlers

import (
	"net/http"
	"strconv"

	"crm-app/backend/models"

	"github.com/gin-gonic/gin"
)

// CRMTargetHandler handles requests for sales targets
type CRMTargetHandler struct {
	targetRepo models.TargetRepository
}

// NewCRMTargetHandler creates a new target handler
func NewCRMTargetHandler(repos *models.CRMRepositories) *CRMTargetHandler {
	return &CRMTargetHandler{
		targetRepo: repos.TargetRepo,
	}
}

// GetTargets returns all targets with optional filtering
func (h *CRMTargetHandler) GetTargets(c *gin.Context) {
	// Parse query parameters for filtering
	targetType := c.Query("target_type")
	assignedToStr := c.Query("assigned_to")
	teamIDStr := c.Query("team_id")
	activeStr := c.Query("active")

	filters := make(map[string]interface{})

	if targetType != "" {
		filters["target_type"] = targetType
	}

	if assignedToStr != "" {
		assignedTo, err := strconv.Atoi(assignedToStr)
		if err == nil {
			filters["assigned_to"] = assignedTo
		}
	}

	if teamIDStr != "" {
		teamID, err := strconv.Atoi(teamIDStr)
		if err == nil {
			filters["team_id"] = teamID
		}
	}

	if activeStr == "true" {
		filters["active"] = true
	} else if activeStr == "false" {
		filters["active"] = false
	}

	targets, err := h.targetRepo.GetTargets(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch targets"})
		return
	}

	c.JSON(http.StatusOK, targets)
}

// GetTarget returns a target by ID
func (h *CRMTargetHandler) GetTarget(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID"})
		return
	}

	target, err := h.targetRepo.GetTargetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch target"})
		return
	}

	if target == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target not found"})
		return
	}

	c.JSON(http.StatusOK, target)
}

// CreateTarget creates a new target
func (h *CRMTargetHandler) CreateTarget(c *gin.Context) {
	var target models.Target
	if err := c.ShouldBindJSON(&target); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the created_by field to the current user ID
	// In a real app, would get this from the authenticated user
	// userID := 1 // Placeholder
	// target.CreatedBy = userID

	if err := h.targetRepo.CreateTarget(&target); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create target"})
		return
	}

	c.JSON(http.StatusCreated, target)
}

// UpdateTarget updates a target
func (h *CRMTargetHandler) UpdateTarget(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID"})
		return
	}

	// Verify that the target exists
	existingTarget, err := h.targetRepo.GetTargetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch target"})
		return
	}
	if existingTarget == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target not found"})
		return
	}

	var target models.Target
	if err := c.ShouldBindJSON(&target); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure ID matches the URL parameter
	target.ID = id

	// Preserve the created_by field
	// target.CreatedBy = existingTarget.CreatedBy

	if err := h.targetRepo.UpdateTarget(&target); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update target"})
		return
	}

	c.JSON(http.StatusOK, target)
}

// DeleteTarget deletes a target
func (h *CRMTargetHandler) DeleteTarget(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID"})
		return
	}

	// Verify that the target exists
	existingTarget, err := h.targetRepo.GetTargetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch target"})
		return
	}
	if existingTarget == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target not found"})
		return
	}

	if err := h.targetRepo.DeleteTarget(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete target"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Target deleted successfully"})
}

// GetTargetProgress returns progress for a specific target
func (h *CRMTargetHandler) GetTargetProgress(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID"})
		return
	}

	progress, err := h.targetRepo.GetTargetProgress(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch target progress"})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// GetAllTargetProgress returns progress for all active targets
func (h *CRMTargetHandler) GetAllTargetProgress(c *gin.Context) {
	progress, err := h.targetRepo.GetAllTargetProgress()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch target progress"})
		return
	}

	c.JSON(http.StatusOK, progress)
}
