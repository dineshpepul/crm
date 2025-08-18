package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"crm-app/backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Init initializes the database connection
func Init() (*gorm.DB, error) {
	// Get database connection details from environment variables
	dbUser := getEnvOrDefault("DB_USER", "root")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "root")
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "8889")
	dbName := getEnvOrDefault("DB_NAME", "crmgo")

	// Create the connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Configure GORM logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Enable color
		},
	)

	// Open database connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// AutoMigrate automatically migrates the schema for all models
func AutoMigrate(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	// Auto migrate all relevant models
	return db.AutoMigrate(
		&models.User{},
		&models.Lead{},
		&models.LeadCustomField{},
		&models.LeadTag{},
		&models.LeadFieldConfig{},
		&models.LeadFormSection{},
		&models.Contact{},
		&models.Deal{},
		&models.Campaign{},
		&models.CampaignLead{},
		&models.CampaignTemplate{},
		&models.Target{},
		&models.NurtureSequence{},
		&models.NurtureStep{},
		&models.NurtureEnrollment{},
		&models.NurtureActivity{},
		&models.CrmFieldData{},
		&models.LeadInput{},
		&models.LeadData{},
	)
}

// Helper function to get environment variable with default fallback
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
