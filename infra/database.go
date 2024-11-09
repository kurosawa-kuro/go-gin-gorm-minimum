package infra

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDB() *gorm.DB {
	// 環境変数から設定を読み込む
	config := loadDBConfig()
	log.Println("config", config)

	var (
		db  *gorm.DB
		err error
	)

	switch config.Environment {
	case "prod":
		db, err = setupPostgres(config)
		log.Printf("Setup postgresql database for %s", config.Environment)

	case "dev":
		db, err = setupPostgres(config)
		log.Printf("Setup postgresql database for %s", config.Environment)

	case "test":
		testConfig := DBConfig{
			Host:     "localhost",
			User:     "postgres",
			Password: "postgres",
			DBName:   "web_app_db_integration_test_go",
			Port:     "5432",
		}
		db, err = setupPostgres(testConfig)
		log.Println("Setup postgresql database for testing")

	default:
		config.Environment = "dev"
		db, err = setupPostgres(config)
		log.Println("Setup postgresql database for development (default)")
	}

	if err != nil {
		panic(fmt.Sprintf("Failed to connect database: %v", err))
	}

	return db
}

type DBConfig struct {
	Environment string
	Host        string
	User        string
	Password    string
	DBName      string
	Port        string
}

func loadDBConfig() DBConfig {
	return DBConfig{
		Environment: getEnvOrDefault("ENV", "dev"),
		Host:        getEnvOrDefault("DB_HOST", "localhost"),
		User:        getEnvOrDefault("DB_USER", "postgres"),
		Password:    getEnvOrDefault("DB_PASSWORD", "postgres"),
		DBName:      getEnvOrDefault("DB_NAME", "web_app_db_integration_go"),
		Port:        getEnvOrDefault("DB_PORT", "5432"),
	}
}

func setupPostgres(config DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo",
		config.Host,
		config.User,
		config.Password,
		config.DBName,
		config.Port,
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
