package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User モデル定義（メインのUserモデルと同じ構造）
type User struct {
	ID         uint   `gorm:"primaryKey"`
	Email      string `gorm:"uniqueIndex;not null"`
	Password   string `gorm:"not null"`
	Role       string `gorm:"default:'user'"`
	AvatarPath string
}

func main() {
	// データベース接続
	dsn := "host=localhost user=postgres password=postgres dbname=web_app_db_integration_go port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// ユーザーテーブルの全レコードを取得
	var users []User
	result := db.Find(&users)
	if result.Error != nil {
		log.Fatal("failed to fetch users:", result.Error)
	}

	// 各ユーザーを削除し、IDを表示
	for _, user := range users {
		fmt.Printf("Deleting user ID: %d, Email: %s\n", user.ID, user.Email)
		if err := db.Delete(&user).Error; err != nil {
			log.Printf("Error deleting user %d: %v\n", user.ID, err)
			continue
		}
	}

	// 削除した総数を表示
	fmt.Printf("\nTotal users deleted: %d\n", len(users))

	// シーケンスをリセット
	if err := db.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Error resetting sequence: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("User ID sequence reset to 1")
}
