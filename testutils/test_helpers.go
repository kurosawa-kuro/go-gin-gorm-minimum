package testutils

import (
	"go-gin-gorm-minimum/handlers"
	"go-gin-gorm-minimum/infra"
	"go-gin-gorm-minimum/models"
	"go-gin-gorm-minimum/services"
	"go-gin-gorm-minimum/utils"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var TestDB *gorm.DB

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

// SetupTestRouter は共通のテスト用Ginルーターをセットアップします
func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.Default()
}

// SetupUserHandler はUserHandlerとその依存関係をセットアップします
func SetupUserHandler() (*gin.Engine, *handlers.UserHandler) {
	if err := CleanupDatabase(TestDB); err != nil {
		panic(err)
	}

	userService := services.NewUserService(TestDB)
	userHandler := handlers.NewUserHandler(userService)
	r := SetupTestRouter()

	return r, userHandler
}

// SetupMicropostHandler はMicropostHandlerとその依存関係をセットアップします
func SetupMicropostHandler() (*gin.Engine, *handlers.MicropostHandler, *services.AuthService) {
	if err := CleanupDatabase(TestDB); err != nil {
		panic(err)
	}

	micropostService := services.NewMicropostService(TestDB)
	micropostHandler := handlers.NewMicropostHandler(micropostService)
	authService := services.NewAuthService(TestDB)
	r := SetupTestRouter()

	return r, micropostHandler, authService
}

// SetupAuthHandler はAuthHandlerとその依存関係をセットアップします
func SetupAuthHandler() (*gin.Engine, *handlers.AuthHandler) {
	if err := CleanupDatabase(TestDB); err != nil {
		panic(err)
	}

	userService := services.NewUserService(TestDB)
	authService := services.NewAuthService(TestDB)
	authHandler := handlers.NewAuthHandler(authService, userService)
	r := SetupTestRouter()

	return r, authHandler
}

// CreateTestUser creates a test user and returns the user object and a JWT token
func CreateTestUser() (*models.User, string, error) {
	user := &models.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	if err := TestDB.Create(user).Error; err != nil {
		return nil, "", err
	}

	// utils\jwt.goのGenerateJWTToken関数を使ってJWTトークンを生成
	token, err := utils.GenerateJWTToken(*user) // utilsパッケージのGenerateJWTToken関数を呼び出す
	return user, token, err
}
