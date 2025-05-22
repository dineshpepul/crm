package routes

import (
	"crm-app/backend/handlers"
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

	// CRM API group
	crm := r.Group("/api/crm")

	// Dashboard routes
	dashboard := crm.Group("/dashboard")
	{
		dashboard.GET("/summary", dashboardHandler.GetDashboardSummary)
		dashboard.GET("/leads-by-source", dashboardHandler.GetLeadsBySource)
		dashboard.GET("/leads-by-status", dashboardHandler.GetLeadsByStatus)
		dashboard.GET("/revenue-by-month", dashboardHandler.GetRevenueByMonth)
		dashboard.GET("/sales-forecast", dashboardHandler.GetSalesForecast)
		dashboard.GET("/top-deals", dashboardHandler.GetTopDeals)
		dashboard.GET("/recent-leads", dashboardHandler.GetRecentLeads)
		dashboard.GET("/target-progress", dashboardHandler.GetTargetProgress)
	}

	// Lead routes
	leads := crm.Group("/leads")
	{
		leads.GET("", leadHandler.GetLeads)
		leads.POST("", leadHandler.CreateLead)
		leads.GET("/:id", leadHandler.GetLead)
		leads.PUT("/:id", leadHandler.UpdateLead)
		leads.DELETE("/:id", leadHandler.DeleteLead)

		// Lead qualification routes
		leads.PUT("/:id/qualify", leadHandler.QualifyLead)
		leads.PUT("/:id/disqualify", leadHandler.DisqualifyLead)

		// Lead assignment routes
		leads.PUT("/:id/assign", leadHandler.AssignLead)

		// Bulk operations
		leads.POST("/import", leadHandler.BulkImportLeads)
		leads.GET("/export", leadHandler.ExportLeads)
	}

	// Lead field configuration routes
	leadFields := crm.Group("/lead-fields")
	{
		// Field configurations
		leadFields.GET("", leadFieldsHandler.GetAllFieldConfigs)
		leadFields.GET("/visible", leadFieldsHandler.GetVisibleFieldConfigs)
		leadFields.GET("/required", leadFieldsHandler.GetRequiredFieldConfigs)
		leadFields.GET("/section/:section", leadFieldsHandler.GetFieldConfigsBySection)
		leadFields.GET("/:id", leadFieldsHandler.GetFieldConfig)
		leadFields.POST("", leadFieldsHandler.CreateFieldConfig)
		leadFields.PUT("/:id", leadFieldsHandler.UpdateFieldConfig)
		leadFields.DELETE("/:id", leadFieldsHandler.DeleteFieldConfig)
		leadFields.POST("/reorder", leadFieldsHandler.ReorderFormFields)

		// Form sections
		leadFields.GET("/sections", leadFieldsHandler.GetAllFormSections)
		leadFields.GET("/sections/visible", leadFieldsHandler.GetVisibleFormSections)
		leadFields.POST("/sections", leadFieldsHandler.CreateFormSection)
		leadFields.PUT("/sections/:id", leadFieldsHandler.UpdateFormSection)
		leadFields.DELETE("/sections/:id", leadFieldsHandler.DeleteFormSection)
		leadFields.POST("/sections/reorder", leadFieldsHandler.ReorderFormSections)

		// Complete form structure
		leadFields.GET("/form-structure", leadFieldsHandler.GetFormStructure)
	}

	// Deal routes
	deals := crm.Group("/deals")
	{
		deals.GET("", dealHandler.GetDeals)
		deals.POST("", dealHandler.CreateDeal)
		deals.GET("/:id", dealHandler.GetDeal)
		deals.PUT("/:id", dealHandler.UpdateDeal)
		deals.DELETE("/:id", dealHandler.DeleteDeal)

		// Deal-specific routes
		deals.PUT("/:id/stage", dealHandler.UpdateDealStage)
		deals.GET("/lead/:lead_id", dealHandler.GetDealsByLead)
		deals.GET("/pipeline", dealHandler.GetDealPipeline)
	}

	// Contact routes
	contacts := crm.Group("/contacts")
	{
		contacts.GET("", contactHandler.GetContacts)
		contacts.POST("", contactHandler.CreateContact)
		contacts.GET("/:id", contactHandler.GetContact)
		contacts.PUT("/:id", contactHandler.UpdateContact)
		contacts.DELETE("/:id", contactHandler.DeleteContact)

		// Contact-specific routes
		contacts.GET("/search", contactHandler.SearchContacts)
		contacts.GET("/lead/:lead_id", contactHandler.GetContactsByLead)
	}

	// Nurture routes
	nurture := crm.Group("/nurture")
	{
		// Campaign routes
		campaigns := nurture.Group("/campaigns")
		{
			campaigns.GET("", nurtureHandler.GetCampaigns)
			campaigns.POST("", nurtureHandler.CreateCampaign)
			campaigns.GET("/:id", nurtureHandler.GetCampaign)
			campaigns.PUT("/:id", nurtureHandler.UpdateCampaign)
			campaigns.DELETE("/:id", nurtureHandler.DeleteCampaign)

			// Campaign-specific routes
			campaigns.GET("/:id/stats", nurtureHandler.GetCampaignStats)
			campaigns.GET("/:id/leads", nurtureHandler.GetCampaignLeads)
			campaigns.POST("/:id/leads", nurtureHandler.AddLeadsToCampaign)
			campaigns.DELETE("/:id/leads", nurtureHandler.RemoveLeadsFromCampaign)
		}

		// Template routes
		templates := nurture.Group("/templates")
		{
			templates.GET("", nurtureHandler.GetTemplates)
			templates.POST("", nurtureHandler.CreateTemplate)
			templates.GET("/:id", nurtureHandler.GetTemplate)
			templates.PUT("/:id", nurtureHandler.UpdateTemplate)
			templates.DELETE("/:id", nurtureHandler.DeleteTemplate)
		}
	}

	// Analytics routes
	analytics := crm.Group("/analytics")
	{
		analytics.GET("/leads", analyticsHandler.GetLeadAnalytics)
		analytics.GET("/deals", analyticsHandler.GetDealAnalytics)
		analytics.GET("/sales-activity", analyticsHandler.GetSalesActivityAnalytics)
		analytics.GET("/performance", analyticsHandler.GetPerformanceAnalytics)
		analytics.GET("/funnel", analyticsHandler.GetFunnelAnalytics)
		analytics.GET("/targets", analyticsHandler.GetTargetAnalytics)
		analytics.GET("/dashboard", analyticsHandler.GetDashboardAnalytics)
		analytics.GET("/conversion", analyticsHandler.GetConversionAnalytics)
	}

	// Target routes
	targets := crm.Group("/targets")
	{
		targets.GET("", targetHandler.GetTargets)
		targets.POST("", targetHandler.CreateTarget)
		targets.GET("/:id", targetHandler.GetTarget)
		targets.PUT("/:id", targetHandler.UpdateTarget)
		targets.DELETE("/:id", targetHandler.DeleteTarget)

		// Target progress routes
		targets.GET("/:id/progress", targetHandler.GetTargetProgress)
		targets.GET("/progress", targetHandler.GetAllTargetProgress)
	}
}
