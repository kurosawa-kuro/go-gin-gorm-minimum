package handlers

import (
	"go-gin-gorm-minimum/models"

	"gorm.io/gorm"
)

type DatabaseOperations struct {
	db *gorm.DB
}

func NewDatabaseOperations(db *gorm.DB) *DatabaseOperations {
	return &DatabaseOperations{db: db}
}

func (ops *DatabaseOperations) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := ops.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
