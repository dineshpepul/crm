package handlers

import (
	"net/http"
	"strconv"

	"crm-app/backend/models"

	"github.com/gin-gonic/gin"
)

// CRMContactHandler handles requests for contact management
type CRMContactHandler struct {
	contactRepo models.ContactRepository
	leadRepo    models.LeadRepository
}

// NewCRMContactHandler creates a new contact handler
func NewCRMContactHandler(repos *models.CRMRepositories) *CRMContactHandler {
	return &CRMContactHandler{
		contactRepo: repos.ContactRepo,
		leadRepo:    repos.LeadRepo,
	}
}

// GetContacts returns all contacts with pagination
func (h *CRMContactHandler) GetContacts(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}

	contacts, err := h.contactRepo.List(offset, limit, companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contacts"})
		return
	}

	c.JSON(http.StatusOK, contacts)
}

// GetContact returns a contact by ID
func (h *CRMContactHandler) GetContact(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	contact, err := h.contactRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contact"})
		return
	}

	if contact == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}

	c.JSON(http.StatusOK, contact)
}

// CreateContact creates a new contact
func (h *CRMContactHandler) CreateContact(c *gin.Context) {
	var contact models.Contact
	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If lead_id is provided, verify that the lead exists
	if contact.LeadID != nil {
		lead, err := h.leadRepo.FindByID(*contact.LeadID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify lead"})
			return
		}
		if lead == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Lead does not exist"})
			return
		}
	}

	if err := h.contactRepo.Create(&contact); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create contact"})
		return
	}

	c.JSON(http.StatusCreated, contact)
}

// UpdateContact updates a contact
func (h *CRMContactHandler) UpdateContact(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	// Verify that the contact exists
	existingContact, err := h.contactRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contact"})
		return
	}
	if existingContact == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}

	var contact models.Contact
	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure ID matches the URL parameter
	contact.ID = id

	// If lead_id is provided, verify that the lead exists
	if contact.LeadID != nil {
		lead, err := h.leadRepo.FindByID(*contact.LeadID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify lead"})
			return
		}
		if lead == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Lead does not exist"})
			return
		}
	}

	if err := h.contactRepo.Update(&contact); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update contact"})
		return
	}

	c.JSON(http.StatusOK, contact)
}

// DeleteContact deletes a contact
func (h *CRMContactHandler) DeleteContact(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	// Verify that the contact exists
	existingContact, err := h.contactRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contact"})
		return
	}
	if existingContact == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}

	if err := h.contactRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete contact"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact deleted successfully"})
}

// GetContactsByLead returns contacts for a specific lead
func (h *CRMContactHandler) GetContactsByLead(c *gin.Context) {
	leadIDStr := c.Param("lead_id")
	leadID, err := strconv.Atoi(leadIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lead ID"})
		return
	}

	contacts, err := h.contactRepo.FindByLead(leadID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contacts"})
		return
	}

	c.JSON(http.StatusOK, contacts)
}

// SearchContacts searches contacts by name or email
func (h *CRMContactHandler) SearchContacts(c *gin.Context) {
	query := c.Query("q")
	companyIdStr := c.Query("companyId")
	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyId"})
		return
	}
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	contacts, err := h.contactRepo.Search(query, companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search contacts"})
		return
	}

	c.JSON(http.StatusOK, contacts)
}
