package handlers_test

import (
	"bytes"
	"encoding/json"
	"go-gin-gorm-minimum/handlers"
	"go-gin-gorm-minimum/infra"
	"go-gin-gorm-minimum/models"
	"go-gin-gorm-minimum/services"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	db          *gorm.DB
	authHandler *handlers.AuthHandler
)

func TestMain(m *testing.M) {
	// テスト用の環境変数設定
	os.Setenv("ENV", "test")

	// テスト用DBのセットアップ
	db = infra.SetupDB()

	// テスト実行
	code := m.Run()

	// クリーンアップ
	sql, err := db.DB()
	if err == nil {
		sql.Close()
	}

	os.Exit(code)
}

func cleanupDatabase(db *gorm.DB) error {
	// 外部キー制約を一時的に無効化
	db.Exec("SET CONSTRAINTS ALL DEFERRED")

	// テーブルのクリーンアップ（順序を考慮）
	if err := db.Exec("DELETE FROM microposts").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM users").Error; err != nil {
		return err
	}

	// 外部キー制約を再度有効化
	db.Exec("SET CONSTRAINTS ALL IMMEDIATE")
	return nil
}

func setupTest() (*gin.Engine, *handlers.AuthHandler) {
	// テスト前にDBをクリーン
	if err := cleanupDatabase(db); err != nil {
		panic(err)
	}

	// 必要なサービスの初期化
	userService := services.NewUserService(db)
	authService := services.NewAuthService(db)
	authHandler := handlers.NewAuthHandler(authService, userService)

	// Ginルーターの設定
	gin.SetMode(gin.TestMode) // テストモードに設定
	r := gin.Default()
	return r, authHandler
}

func TestLoginUser(t *testing.T) {
	r, authHandler := setupTest()
	r.POST("/auth/login", authHandler.LoginUser)

	// テストケース
	tests := []struct {
		name           string
		setupUser      *models.User
		loginRequest   models.LoginRequest
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Successful Login",
			setupUser: &models.User{
				Email:    "test1@example.com",
				Password: "password123",
			},
			loginRequest: models.LoginRequest{
				Email:    "test1@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.LoginResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.Token)
			},
		},
		{
			name: "Invalid Credentials",
			setupUser: &models.User{
				Email:    "test2@example.com",
				Password: "password123",
			},
			loginRequest: models.LoginRequest{
				Email:    "test2@example.com",
				Password: "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		{
			name:      "User Not Found",
			setupUser: nil,
			loginRequest: models.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト前にDBをクリーン
			err := cleanupDatabase(db)
			assert.NoError(t, err)

			// テストユーザーのセットアップ
			if tt.setupUser != nil {
				err := services.NewAuthService(db).SignUp(tt.setupUser)
				assert.NoError(t, err)
			}

			// リクエストの作成
			loginJSON, err := json.Marshal(tt.loginRequest)
			assert.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(loginJSON))
			req.Header.Set("Content-Type", "application/json")

			// レスポンスの記録
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			// アサーション
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}
