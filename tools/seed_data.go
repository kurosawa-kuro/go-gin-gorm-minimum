package main

// tools\delete_data.go

import (
	"fmt"
	"go-gin-gorm-minimum/models"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SeedData() {
	// データベース接続
	dsn := "host=localhost user=postgres password=postgres dbname=web_app_db_integration_go port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// ユーザーとマイクロポストのデータを挿入 ユーザーのパスワードは""でハッシュ化して
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("failed to hash password:", err)
	}
	SeedUsersAndMicroposts(db, hashedPassword)

	fmt.Println("\nAll data inserted successfully")
}

func SeedUsersAndMicroposts(db *gorm.DB, hashedPassword []byte) {
	// サンプルユーザーの作成
	users := []models.User{
		{
			Email:    "user1@example.com",
			Password: string(hashedPassword),
		},
		{
			Email:    "user2@example.com",
			Password: string(hashedPassword),
		},
	}

	// ユーザーの保存
	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			log.Printf("Error creating user %s: %v\n", user.Email, err)
			continue
		}
		fmt.Printf("Created user: %s\n", user.Email)

		// 各ユーザーのマイクロポストを作成
		microposts := []models.Micropost{
			{
				Title: fmt.Sprintf("First post by %s", user.Email),
			},
			{
				Title: fmt.Sprintf("Second post by %s", user.Email),
			},
		}

		// マイクロポストの保存
		for _, post := range microposts {
			if err := db.Create(&post).Error; err != nil {
				log.Printf("Error creating micropost for user %d: %v\n", user.ID, err)
				continue
			}
			fmt.Printf("Created micropost: %s\n", post.Title)
		}
	}
}
