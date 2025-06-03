package routes

import (
	"crm-app/backend/handlers"
	"crm-app/backend/middleware"
	"crm-app/backend/models"

	"github.com/gin-gonic/gin"
)

// SetupCRMRoutes sets up all CRM-related routes
func SetupCRMRoutes(r *gin.Engine, repos *models.CRMRepositories) {
	// Initialize services
	// Services are initialized in the handlers constructors

	// Initialize handlers
	dashboardHandler := handlers.NewCRMDashboardHandler(repos)
	leadHandler := handlers.NewCRMLeadHandler(repos)
	dealHandler := handlers.NewCRMDealHandler(repos)
	contactHandler := handlers.NewCRMContactHandler(repos)
	nurtureHandler := handlers.NewCRMNurtureHandler(repos)
	analyticsHandler := handlers.NewCRMAnalyticsHandler(repos)
	targetHandler := handlers.NewCRMTargetHandler(repos)
	leadFieldsHandler := handlers.NewCRMLeadFieldsHandler(repos)
	LeadScoreHandler := handlers.NewScoreLeadHandler(repos)

	// CRM API group
	crm := r.Group("/api/crm")

	// Dashboard routes
	dashboard := crm.Group("/dashboard")
	{
		dashboard.GET("/summary", middleware.JwtAuthMiddleware(), dashboardHandler.GetDashboardSummary)
		dashboard.GET("/leads-by-source", middleware.JwtAuthMiddleware(), dashboardHandler.GetLeadsBySource)
		dashboard.GET("/leads-by-status", middleware.JwtAuthMiddleware(), dashboardHandler.GetLeadsByStatus)
		dashboard.GET("/revenue-by-month", middleware.JwtAuthMiddleware(), dashboardHandler.GetRevenueByMonth)
		dashboard.GET("/sales-forecast", middleware.JwtAuthMiddleware(), dashboardHandler.GetSalesForecast)
		dashboard.GET("/top-deals", middleware.JwtAuthMiddleware(), dashboardHandler.GetTopDeals)
		dashboard.GET("/recent-leads", middleware.JwtAuthMiddleware(), dashboardHandler.GetRecentLeads)
		dashboard.GET("/target-progress", middleware.JwtAuthMiddleware(), dashboardHandler.GetTargetProgress)
	}
	// Lead routes
	leads := crm.Group("/leads")
	{
		leads.GET("", middleware.JwtAuthMiddleware(), leadHandler.GetLeads)
		leads.POST("", middleware.JwtAuthMiddleware(), leadHandler.CreateLead)
		leads.GET("/:id", middleware.JwtAuthMiddleware(), leadHandler.GetLead)
		leads.PUT("/:id", middleware.JwtAuthMiddleware(), leadHandler.UpdateLead)
		leads.DELETE("/:id", middleware.JwtAuthMiddleware(), leadHandler.DeleteLead)

		// Lead qualification routes
		leads.PUT("/:id/qualify", middleware.JwtAuthMiddleware(), leadHandler.QualifyLead)
		leads.PUT("/:id/disqualify", middleware.JwtAuthMiddleware(), leadHandler.DisqualifyLead)

		// Lead assignment routes
		leads.PUT("/:id/assign", middleware.JwtAuthMiddleware(), leadHandler.AssignLead)
		leads.PUT("/updateScore", middleware.JwtAuthMiddleware(), LeadScoreHandler.UpdateScore)

		// Bulk operations
		leads.POST("/import", middleware.JwtAuthMiddleware(), leadHandler.BulkImportLeads)
		leads.GET("/export", middleware.JwtAuthMiddleware(), leadHandler.ExportLeads)
	}

	// Lead field configuration routes
	leadFields := crm.Group("/lead-fields")
	{
		// Field configurations
		leadFields.GET("", middleware.JwtAuthMiddleware(), leadFieldsHandler.GetAllFieldConfigs)
		leadFields.GET("/visible", middleware.JwtAuthMiddleware(), leadFieldsHandler.GetVisibleFieldConfigs)
		leadFields.GET("/required", middleware.JwtAuthMiddleware(), leadFieldsHandler.GetRequiredFieldConfigs)
		leadFields.GET("/section/:section", leadFieldsHandler.GetFieldConfigsBySection)
		leadFields.GET("/:id", middleware.JwtAuthMiddleware(), leadFieldsHandler.GetFieldConfig)
		leadFields.POST("", middleware.JwtAuthMiddleware(), leadFieldsHandler.CreateFieldConfig)
		leadFields.PUT("/:id", middleware.JwtAuthMiddleware(), leadFieldsHandler.UpdateFieldConfig)
		leadFields.DELETE("/:id", middleware.JwtAuthMiddleware(), leadFieldsHandler.DeleteFieldConfig)
		leadFields.POST("/reorder", middleware.JwtAuthMiddleware(), leadFieldsHandler.ReorderFormFields)

		// Form sections
		leadFields.GET("/sections", middleware.JwtAuthMiddleware(), leadFieldsHandler.GetAllFormSections)
		leadFields.GET("/sections/visible", middleware.JwtAuthMiddleware(), leadFieldsHandler.GetVisibleFormSections)
		leadFields.POST("/sections", middleware.JwtAuthMiddleware(), leadFieldsHandler.CreateFormSection)
		leadFields.PUT("/sections/:id", middleware.JwtAuthMiddleware(), leadFieldsHandler.UpdateFormSection)
		leadFields.DELETE("/sections/:id", middleware.JwtAuthMiddleware(), leadFieldsHandler.DeleteFormSection)
		leadFields.POST("/sections/reorder", middleware.JwtAuthMiddleware(), leadFieldsHandler.ReorderFormSections)

		// Complete form structure
		leadFields.GET("/form-structure", middleware.JwtAuthMiddleware(), leadFieldsHandler.GetFormStructure)
	}

	// Deal routes
	deals := crm.Group("/deals")
	{
		deals.GET("", middleware.JwtAuthMiddleware(), dealHandler.GetDeals)
		deals.POST("", middleware.JwtAuthMiddleware(), dealHandler.CreateDeal)
		deals.GET("/:id", middleware.JwtAuthMiddleware(), dealHandler.GetDeal)
		deals.PUT("/:id", middleware.JwtAuthMiddleware(), dealHandler.UpdateDeal)
		deals.DELETE("/:id", middleware.JwtAuthMiddleware(), dealHandler.DeleteDeal)

		// Deal-specific routes
		deals.PUT("/:id/stage", middleware.JwtAuthMiddleware(), dealHandler.UpdateDealStage)
		deals.GET("/lead/:lead_id", middleware.JwtAuthMiddleware(), dealHandler.GetDealsByLead)
		deals.GET("/pipeline", middleware.JwtAuthMiddleware(), dealHandler.GetDealPipeline)
	}

	// Contact routes
	contacts := crm.Group("/contacts")
	{
		contacts.GET("", middleware.JwtAuthMiddleware(), contactHandler.GetContacts)
		contacts.POST("", middleware.JwtAuthMiddleware(), contactHandler.CreateContact)
		contacts.GET("/:id", middleware.JwtAuthMiddleware(), contactHandler.GetContact)
		contacts.PUT("/:id", middleware.JwtAuthMiddleware(), contactHandler.UpdateContact)
		contacts.DELETE("/:id", middleware.JwtAuthMiddleware(), contactHandler.DeleteContact)

		// Contact-specific routes
		contacts.GET("/search", middleware.JwtAuthMiddleware(), contactHandler.SearchContacts)
		contacts.GET("/lead/:lead_id", middleware.JwtAuthMiddleware(), contactHandler.GetContactsByLead)
	}

	// Nurture routes
	nurture := crm.Group("/nurture")
	{
		// Campaign routes
		campaigns := nurture.Group("/campaigns")
		{
			campaigns.GET("", middleware.JwtAuthMiddleware(), nurtureHandler.GetCampaigns)
			campaigns.POST("", middleware.JwtAuthMiddleware(), nurtureHandler.CreateCampaign)
			campaigns.GET("/:id", middleware.JwtAuthMiddleware(), nurtureHandler.GetCampaign)
			campaigns.PUT("/:id", middleware.JwtAuthMiddleware(), nurtureHandler.UpdateCampaign)
			campaigns.DELETE("/:id", middleware.JwtAuthMiddleware(), nurtureHandler.DeleteCampaign)

			// Campaign-specific routes
			campaigns.GET("/:id/stats", middleware.JwtAuthMiddleware(), nurtureHandler.GetCampaignStats)
			campaigns.GET("/:id/leads", middleware.JwtAuthMiddleware(), nurtureHandler.GetCampaignLeads)
			campaigns.POST("/:id/leads", middleware.JwtAuthMiddleware(), nurtureHandler.AddLeadsToCampaign)
			campaigns.DELETE("/:id/leads", middleware.JwtAuthMiddleware(), nurtureHandler.RemoveLeadsFromCampaign)
		}

		// Template routes
		templates := nurture.Group("/templates")
		{
			templates.GET("", middleware.JwtAuthMiddleware(), nurtureHandler.GetTemplates)
			templates.POST("", middleware.JwtAuthMiddleware(), nurtureHandler.CreateTemplate)
			templates.GET("/:id", middleware.JwtAuthMiddleware(), nurtureHandler.GetTemplate)
			templates.PUT("/:id", middleware.JwtAuthMiddleware(), nurtureHandler.UpdateTemplate)
			templates.DELETE("/:id", middleware.JwtAuthMiddleware(), nurtureHandler.DeleteTemplate)
		}
	}

	// Analytics routes
	analytics := crm.Group("/analytics")
	{
		analytics.GET("/leads", middleware.JwtAuthMiddleware(), analyticsHandler.GetLeadAnalytics)
		analytics.GET("/deals", middleware.JwtAuthMiddleware(), analyticsHandler.GetDealAnalytics)
		analytics.GET("/sales-activity", middleware.JwtAuthMiddleware(), analyticsHandler.GetSalesActivityAnalytics)
		analytics.GET("/performance", middleware.JwtAuthMiddleware(), analyticsHandler.GetPerformanceAnalytics)
		analytics.GET("/funnel", middleware.JwtAuthMiddleware(), analyticsHandler.GetFunnelAnalytics)
		analytics.GET("/targets", middleware.JwtAuthMiddleware(), analyticsHandler.GetTargetAnalytics)
		analytics.GET("/dashboard", middleware.JwtAuthMiddleware(), analyticsHandler.GetDashboardAnalytics)
		analytics.GET("/conversion", middleware.JwtAuthMiddleware(), analyticsHandler.GetConversionAnalytics)
	}

	// Target routes
	targets := crm.Group("/targets")
	{
		targets.GET("", middleware.JwtAuthMiddleware(), targetHandler.GetTargets)
		targets.POST("", middleware.JwtAuthMiddleware(), targetHandler.CreateTarget)
		targets.GET("/:id", middleware.JwtAuthMiddleware(), targetHandler.GetTarget)
		targets.PUT("/:id", middleware.JwtAuthMiddleware(), targetHandler.UpdateTarget)
		targets.DELETE("/:id", middleware.JwtAuthMiddleware(), targetHandler.DeleteTarget)

		// Target progress routes
		targets.GET("/:id/progress", middleware.JwtAuthMiddleware(), targetHandler.GetTargetProgress)
		targets.GET("/progress", middleware.JwtAuthMiddleware(), targetHandler.GetAllTargetProgress)
	}
}
