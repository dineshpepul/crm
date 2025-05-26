package routes

import (
	"crm-app/backend/handlers"
	"crm-app/backend/models"
	"crm-app/backend/services"

	"crm-app/backend/middleware"

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
		leads.GET("/", middleware.JwtAuthMiddleware(), leadHandler.GetLeads)
		leads.POST("/", middleware.JwtAuthMiddleware(), leadHandler.CreateLead)
		leads.GET("/:id", middleware.JwtAuthMiddleware(), leadHandler.GetLead)
		leads.PUT("/:id", middleware.JwtAuthMiddleware(), leadHandler.UpdateLead)
		leads.DELETE("/:id", middleware.JwtAuthMiddleware(), leadHandler.DeleteLead)

		// Lead qualification routes
		leads.POST("/:id/qualify", middleware.JwtAuthMiddleware(), leadHandler.QualifyLead)
		leads.POST("/:id/disqualify", middleware.JwtAuthMiddleware(), leadHandler.DisqualifyLead)

		// Lead assignment routes
		leads.POST("/:id/assign", middleware.JwtAuthMiddleware(), leadHandler.AssignLead)

		// Lead bulk operations
		leads.POST("/import", middleware.JwtAuthMiddleware(), leadHandler.BulkImportLeads)
		leads.GET("/export", middleware.JwtAuthMiddleware(), leadHandler.ExportLeads)
	}

	// Lead field configuration routes
	fieldConfigs := router.Group("/lead-fields")
	{
		fieldConfigs.GET("/", middleware.JwtAuthMiddleware(), leadHandler.GetAllFieldConfigs)
		fieldConfigs.GET("/visible", middleware.JwtAuthMiddleware(), leadHandler.GetVisibleFieldConfigs)
		fieldConfigs.GET("/required", middleware.JwtAuthMiddleware(), leadHandler.GetRequiredFieldConfigs)
		fieldConfigs.GET("/section/:section", middleware.JwtAuthMiddleware(), leadHandler.GetFieldConfigsBySection)
		fieldConfigs.POST("/", middleware.JwtAuthMiddleware(), leadHandler.CreateFieldConfig)
		fieldConfigs.PUT("/:id", middleware.JwtAuthMiddleware(), leadHandler.UpdateFieldConfig)
		fieldConfigs.DELETE("/:id", middleware.JwtAuthMiddleware(), leadHandler.DeleteFieldConfig)
		fieldConfigs.POST("/reorder", middleware.JwtAuthMiddleware(), leadHandler.ReorderFormFields)

		// Form sections
		fieldConfigs.GET("/sections", middleware.JwtAuthMiddleware(), leadHandler.GetAllFormSections)
		fieldConfigs.GET("/sections/visible", middleware.JwtAuthMiddleware(), leadHandler.GetVisibleFormSections)
		fieldConfigs.POST("/sections", middleware.JwtAuthMiddleware(), leadHandler.CreateFormSection)
		fieldConfigs.PUT("/sections/:id", middleware.JwtAuthMiddleware(), leadHandler.UpdateFormSection)
		fieldConfigs.DELETE("/sections/:id", middleware.JwtAuthMiddleware(), leadHandler.DeleteFormSection)
		fieldConfigs.POST("/sections/reorder", middleware.JwtAuthMiddleware(), leadHandler.ReorderFormSections)
	}
}
