package routes

import (
	"crm-app/backend/handlers"
	"crm-app/backend/models"
	"crm-app/backend/services"

	"github.com/gin-gonic/gin"
)

func SetupLeadRoutes(router *gin.RouterGroup, repos *models.Repositories) {
	leadService := services.NewLeadService(repos)
	leadHandler := handlers.NewLeadHandler(leadService)

	// Lead management - using a different path prefix to avoid conflicts
	leads := router.Group("/lead-management")
	{
		leads.GET("/", leadHandler.GetLeads)
		leads.POST("/", leadHandler.CreateLead)
		leads.GET("/:id", leadHandler.GetLead)
		leads.PUT("/:id", leadHandler.UpdateLead)
		leads.DELETE("/:id", leadHandler.DeleteLead)

		// Special actions
		leads.POST("/:id/qualify", leadHandler.QualifyLead)
		leads.POST("/:id/disqualify", leadHandler.DisqualifyLead)
		leads.POST("/:id/assign", leadHandler.AssignLead)
	}

	// Lead field configuration
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
	}
}
