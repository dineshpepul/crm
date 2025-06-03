package main

import (
	"crm-app/backend/db"
	"crm-app/backend/models"
	"crm-app/backend/repositories"
	"crm-app/backend/routes"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database connection
	database, err := db.Init()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate database schema
	// if err := db.AutoMigrate(database); err != nil {
	// 	log.Fatalf("Failed to migrate database schema: %v", err)
	// }

	// Initialize repositories
	repos := repositories.NewRepositoriesInit(database)

	// Initialize router
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Initialize route groups
	api := r.Group("/api")

	// Setup CRM routes with conversion from Repositories to CRMRepositories if needed
	crmRepos := &models.CRMRepositories{
		LeadRepo:            repos.LeadRepo,
		LeadFieldConfigRepo: repos.LeadFieldConfigRepo,
		ContactRepo:         repos.ContactRepo,
		DealRepo:            repos.DealRepo,
		CampaignRepo:        repos.CampaignRepo,
		DashboardRepo:       repos.DashboardRepo,
		AnalyticsRepo:       repos.AnalyticsRepo,
		TargetRepo:          repos.TargetRepo,
		NurtureRepo:         repos.NurtureRepo,
		UserRepo:            repos.UserRepo,
		LeadScoreType:       repos.ScoreRepo,
	}
	routes.SetupCRMRoutes(r, crmRepos)

	// Setup Lead Capture routes - these will be at /api/leads/...
	routes.SetupLeadCaptureRoutes(api, repos)

	// Setup Lead routes - these will be at /api/lead-management/...
	// routes.SetupLeadRoutes(api, repos)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Server running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
