package infra

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Environment string
	Host        string
	User        string
	Password    string
	DBName      string
	Port        string
}

func SetupDB() *gorm.DB {
	config := loadDBConfig()
	log.Printf("★★★ Database Config ★★★ Environment=%s, Host=%s, DBName=%s",
		config.Environment, config.Host, config.DBName)

	db, err := setupPostgres(config)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect database: %v", err))
	}

	return db
}

func loadDBConfig() DBConfig {
	env := getEnvOrDefault("ENV", "dev")
	loadEnvFile(env)

	return DBConfig{
		Environment: env,
		Host:        getEnvOrDefault("DB_HOST", "localhost"),
		User:        getEnvOrDefault("DB_USER", "postgres"),
		Password:    getEnvOrDefault("DB_PASSWORD", "postgres"),
		DBName:      getEnvOrDefault("DB_NAME", "web_app_db_integration_go"),
		Port:        getEnvOrDefault("DB_PORT", "5432"),
	}
}

func loadEnvFile(env string) {
	envFile := fmt.Sprintf(".env.%s", env)
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Warning: Error loading %s file: %v", envFile, err)
	}
}

func setupPostgres(config DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo",
		config.Host, config.User, config.Password, config.DBName, config.Port,
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
