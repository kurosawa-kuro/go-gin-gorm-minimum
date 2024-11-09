package models

import "time"

// User モデル定義
type User struct {
	ID         uint        `json:"id" gorm:"primaryKey" example:"1"`
	Email      string      `json:"email" gorm:"uniqueIndex;not null" binding:"required,email" example:"user1@example.com"`
	Password   string      `json:"password" gorm:"not null" binding:"required,min=6" example:"password123"`
	Role       string      `json:"role" gorm:"default:'user'" example:"user"`
	AvatarPath string      `json:"avatar_path" example:"/avatars/default.png"`
	Microposts []Micropost `json:"microposts,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt  time.Time   `json:"-" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time   `json:"-" gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

// UserResponse は、パスワードを除外したユーザー情報のレスポンス構造体
type UserResponse struct {
	ID         uint      `json:"id" example:"1"`
	Email      string    `json:"email" example:"user1@example.com"`
	Role       string    `json:"role" example:"user"`
	AvatarPath string    `json:"avatar_path" example:"/avatars/default.png"`
	CreatedAt  time.Time `json:"created_at" example:"2024-11-09T18:00:00+09:00"`
	UpdatedAt  time.Time `json:"updated_at" example:"2024-11-09T18:00:00+09:00"`
}

// LoginResponse は、ログイン時のレスポンス構造体
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	UserResponse
}

// ToResponse は User モデルを UserResponse に変換するヘルパー関数
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:         u.ID,
		Email:      u.Email,
		Role:       u.Role,
		AvatarPath: u.AvatarPath,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

// LoginRequest はログインリクエスト用の構造体
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user1@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
}
