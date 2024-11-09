package main

func main() {
	// 1. まずテーブル構造を作成/更新
	Migrate()
	// 2. 既存データを削除
	DeleteData()
	// 3. 新しいデータを投入
	SeedData()
}

//
