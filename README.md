# go-gin-gorm-minimum

http://localhost:8080/swagger/index.html

Todo
- マイグレーション コード
- 環境変数管理 コード


$env:ENV="dev"; go run main.go

func loadDBConfig() DBConfig {
	// この時点でどの環境変数かで調整する必要がある