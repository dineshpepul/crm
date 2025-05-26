package routes

import (
	"crm-app/backend/handlers"
	"crm-app/backend/middleware"
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
		leads.GET("/", middleware.JwtAuthMiddleware(), leadHandler.GetLeads)
		leads.POST("/", middleware.JwtAuthMiddleware(), leadHandler.CreateLead)
		leads.GET("/:id", middleware.JwtAuthMiddleware(), leadHandler.GetLead)
		leads.PUT("/:id", middleware.JwtAuthMiddleware(), leadHandler.UpdateLead)
		leads.DELETE("/:id", middleware.JwtAuthMiddleware(), leadHandler.DeleteLead)

		// Special actions
		leads.POST("/:id/qualify", middleware.JwtAuthMiddleware(), leadHandler.QualifyLead)
		leads.POST("/:id/disqualify", middleware.JwtAuthMiddleware(), leadHandler.DisqualifyLead)
		leads.POST("/:id/assign", middleware.JwtAuthMiddleware(), leadHandler.AssignLead)
	}

	// Lead field configuration
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
	}
}
