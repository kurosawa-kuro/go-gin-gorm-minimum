package models

import "time"

// MicropostRequest はマイクロポスト作成リクエスト用の構造体
type MicropostRequest struct {
	Title string `json:"title" binding:"required" example:"マイクロポストのタイトル"`
}

// Micropost モデル定義
type Micropost struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" binding:"required" example:"マイクロポストのタイトル"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	User      User      `json:"-" gorm:"foreignKey:UserID;references:ID"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

// MicropostResponse はマイクロポストのレスポンス用の構造体
type MicropostResponse struct {
	ID        uint         `json:"id"`
	Title     string       `json:"title"`
	User      UserResponse `json:"user"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// ToResponse は Micropost モデルを MicropostResponse に変換するヘルパー関数
func (m *Micropost) ToResponse() MicropostResponse {
	return MicropostResponse{
		ID:        m.ID,
		Title:     m.Title,
		User:      m.User.ToResponse(),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
