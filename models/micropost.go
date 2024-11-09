package models

import "time"

// Micropost モデル定義
type Micropost struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" binding:"required" example:"マイクロポストのタイトル"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID;references:ID" swaggerignore:"true"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime"`
}
