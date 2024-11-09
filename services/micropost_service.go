package services

import (
	"go-gin-gorm-minimum/models"

	"gorm.io/gorm"
)

type MicropostService struct {
	db *gorm.DB
}

func NewMicropostService(db *gorm.DB) *MicropostService {
	return &MicropostService{db: db}
}

func (s *MicropostService) Create(micropost *models.Micropost) error {
	return s.db.Create(micropost).Error
}

func (s *MicropostService) GetByID(id string) (*models.Micropost, error) {
	var micropost models.Micropost
	err := s.db.First(&micropost, id).Error
	return &micropost, err
}

func (s *MicropostService) GetAll() ([]models.Micropost, error) {
	var microposts []models.Micropost
	err := s.db.Preload("User").Find(&microposts).Error
	return microposts, err
}
