package handlers_test

import (
	"bytes"
	"encoding/json"
	"go-gin-gorm-minimum/handlers"
	"go-gin-gorm-minimum/infra"
	"go-gin-gorm-minimum/models"
	"go-gin-gorm-minimum/services"
	"go-gin-gorm-minimum/testutils"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Set test environment and initialize DB before running tests
	os.Setenv("ENV", "test")
	testutils.TestDB = infra.SetupDB()
}

func setupTest() (*gin.Engine, *handlers.AuthHandler) {
	// テスト前にDBをクリーン
	if err := testutils.CleanupDatabase(testutils.TestDB); err != nil {
		panic(err)
	}

	// 必要なサービスの初期化
	userService := services.NewUserService(testutils.TestDB)
	authService := services.NewAuthService(testutils.TestDB)
	authHandler := handlers.NewAuthHandler(authService, userService)

	// Ginルーターの設定
	gin.SetMode(gin.TestMode)
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
			err := testutils.CleanupDatabase(testutils.TestDB)
			assert.NoError(t, err)

			// テストユーザーのセットアップ
			if tt.setupUser != nil {
				err := services.NewAuthService(testutils.TestDB).SignUp(tt.setupUser)
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
