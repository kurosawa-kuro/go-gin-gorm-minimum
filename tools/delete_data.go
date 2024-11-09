package main

// tools\delete_data.go

import (
	"fmt"
	"log"

	"go-gin-gorm-minimum/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	DeleteData()
}

func DeleteData() {
	// データベース接続
	dsn := "host=localhost user=postgres password=postgres dbname=web_app_db_integration_go port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// ユーザーテーブルの削除処理
	var users []models.User
	if result := db.Find(&users); result.Error != nil {
		log.Fatal("failed to fetch users:", result.Error)
	}

	for _, user := range users {
		fmt.Printf("Deleting user ID: %d, Email: %s\n", user.ID, user.Email)
		if err := db.Delete(&user).Error; err != nil {
			log.Printf("Error deleting user %d: %v\n", user.ID, err)
		}
	}
	fmt.Printf("\nTotal users deleted: %d\n", len(users))

	// マイクロポストテーブルの削除処理
	var microposts []models.Micropost
	if result := db.Find(&microposts); result.Error != nil {
		log.Fatal("failed to fetch microposts:", result.Error)
	}

	for _, micropost := range microposts {
		fmt.Printf("Deleting micropost ID: %d, Title: %s\n", micropost.ID, micropost.Title)
		if err := db.Delete(&micropost).Error; err != nil {
			log.Printf("Error deleting micropost %d: %v\n", micropost.ID, err)
		}
	}
	fmt.Printf("\nTotal microposts deleted: %d\n", len(microposts))

	// シーケンスのリセット
	sequences := []string{"users_id_seq", "microposts_id_seq"}
	for _, seq := range sequences {
		if err := db.Exec(fmt.Sprintf("ALTER SEQUENCE %s RESTART WITH 1", seq)).Error; err != nil {
			log.Printf("Error resetting sequence %s: %v\n", seq, err)
		} else {
			fmt.Printf("%s reset to 1\n", seq)
		}
	}

	fmt.Println("\nAll data deleted and sequences reset successfully")
}
