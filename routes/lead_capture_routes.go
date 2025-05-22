package routes

import (
	"crm-app/backend/handlers"
	"crm-app/backend/models"
	"crm-app/backend/services"

	"github.com/gin-gonic/gin"
)

// SetupLeadCaptureRoutes sets up all Lead Capture related routes
func SetupLeadCaptureRoutes(router *gin.RouterGroup, repos *models.Repositories) {
	// Initialize services
	leadService := services.NewLeadService(repos)

	// Initialize handlers
	leadHandler := handlers.NewLeadHandler(leadService)

	// Lead routes
	leads := router.Group("/leads")
	{
		leads.GET("/", leadHandler.GetLeads)
		leads.POST("/", leadHandler.CreateLead)
		leads.GET("/:id", leadHandler.GetLead)
		leads.PUT("/:id", leadHandler.UpdateLead)
		leads.DELETE("/:id", leadHandler.DeleteLead)

		// Lead qualification routes
		leads.POST("/:id/qualify", leadHandler.QualifyLead)
		leads.POST("/:id/disqualify", leadHandler.DisqualifyLead)

		// Lead assignment routes
		leads.POST("/:id/assign", leadHandler.AssignLead)

		// Lead bulk operations
		leads.POST("/import", leadHandler.BulkImportLeads)
		leads.GET("/export", leadHandler.ExportLeads)
	}

	// Lead field configuration routes
	fieldConfigs := router.Group("/lead-fields")
	{
		fieldConfigs.GET("/", leadHandler.GetAllFieldConfigs)
		fieldConfigs.GET("/visible", leadHandler.GetVisibleFieldConfigs)
		fieldConfigs.GET("/required", leadHandler.GetRequiredFieldConfigs)
		fieldConfigs.GET("/section/:section", leadHandler.GetFieldConfigsBySection)
		fieldConfigs.POST("/", leadHandler.CreateFieldConfig)
		fieldConfigs.PUT("/:id", leadHandler.UpdateFieldConfig)
		fieldConfigs.DELETE("/:id", leadHandler.DeleteFieldConfig)
		fieldConfigs.POST("/reorder", leadHandler.ReorderFormFields)

		// Form sections
		fieldConfigs.GET("/sections", leadHandler.GetAllFormSections)
		fieldConfigs.GET("/sections/visible", leadHandler.GetVisibleFormSections)
		fieldConfigs.POST("/sections", leadHandler.CreateFormSection)
		fieldConfigs.PUT("/sections/:id", leadHandler.UpdateFormSection)
		fieldConfigs.DELETE("/sections/:id", leadHandler.DeleteFormSection)
		fieldConfigs.POST("/sections/reorder", leadHandler.ReorderFormSections)
	}
}
