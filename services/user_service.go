package services

import (
	"go-gin-gorm-minimum/models"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) FindByEmail(email string) (models.User, error) {
	var user models.User
	err := s.db.Where("email = ?", email).First(&user).Error
	return user, err
}

func (s *UserService) FindByID(id interface{}) (models.User, error) {
	var user models.User
	err := s.db.First(&user, id).Error
	return user, err
}

func (s *UserService) FindAll() ([]models.User, error) {
	var users []models.User
	err := s.db.Find(&users).Error
	return users, err
}

func (s *UserService) UpdateAvatar(userID uint, avatarPath string) (models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return user, err
	}

	user.AvatarPath = avatarPath
	err := s.db.Save(&user).Error
	return user, err
}
