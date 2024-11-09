package testutils

import (
	"go-gin-gorm-minimum/infra"
	"os"
	"testing"

	"gorm.io/gorm"
)

var TestDB *gorm.DB // Changed to exported

func TestMain(m *testing.M) {
	os.Setenv("ENV", "test")
	TestDB = infra.SetupDB()

	code := m.Run()

	sql, err := TestDB.DB()
	if err == nil {
		sql.Close()
	}
	os.Exit(code)
}

func CleanupDatabase(db *gorm.DB) error {
	db.Exec("SET CONSTRAINTS ALL DEFERRED")
	if err := db.Exec("DELETE FROM microposts").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM users").Error; err != nil {
		return err
	}
	db.Exec("SET CONSTRAINTS ALL IMMEDIATE")
	return nil
}
