package utils

import (
	"go-gin-gorm-minimum/models"

	"gorm.io/gorm"
)

// Database connection string
const DSN = "host=localhost user=postgres password=postgres dbname=web_app_db_integration_go port=5432 sslmode=disable"

// FindUserByEmail retrieves a user by email
func FindUserByEmail(db *gorm.DB, email string) (models.User, error) {
	var user models.User
	err := db.Where("email = ?", email).First(&user).Error
	return user, err
}

// FindUserByID retrieves a user by ID
func FindUserByID(db *gorm.DB, id interface{}) (models.User, error) {
	var user models.User
	err := db.First(&user, id).Error
	return user, err
}
