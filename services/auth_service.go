package services

import (
	"go-gin-gorm-minimum/models"
	"go-gin-gorm-minimum/utils"

	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) SignUp(user *models.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return s.db.Create(user).Error
}

func (s *AuthService) Login(email, password string) (*models.LoginResponse, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	if err := utils.CheckPassword(user.Password, password); err != nil {
		return nil, err
	}

	tokenString, err := utils.GenerateJWTToken(user)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token:        tokenString,
		UserResponse: user.ToResponse(),
	}, nil
}

func (s *AuthService) GetUserByID(userID interface{}) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
