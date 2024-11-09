package main

import (
	"fmt"
	"log"

	"go-gin-gorm-minimum/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Migrate() {
	// データベース接続
	dsn := "host=localhost user=postgres password=postgres dbname=web_app_db_integration_go port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// マイグレーションの実行
	fmt.Println("Starting database migration...")

	// 1. まず既存のテーブルをドロップ
	db.Migrator().DropTable(&models.Micropost{}, &models.User{})

	// 2. テーブルを再作成
	err = db.AutoMigrate(
		&models.User{},
		&models.Micropost{},
	)

	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	fmt.Println("Database migration completed successfully")
}
