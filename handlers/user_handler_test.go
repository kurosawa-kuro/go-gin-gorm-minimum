package handlers_test

import (
	"encoding/json"
	"go-gin-gorm-minimum/handlers"
	"go-gin-gorm-minimum/infra"
	"go-gin-gorm-minimum/middlewares"
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
	os.Setenv("ENV", "test")
	testutils.TestDB = infra.SetupDB()
}

func setupUserTest() (*gin.Engine, *handlers.UserHandler) {
	if err := testutils.CleanupDatabase(testutils.TestDB); err != nil {
		panic(err)
	}

	userService := services.NewUserService(testutils.TestDB)
	userHandler := handlers.NewUserHandler(userService)

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// ミドルウェアを設定
	auth := r.Group("/")
	auth.Use(middlewares.AuthMiddleware())
	auth.GET("/users", userHandler.GetUsers)

	return r, userHandler
}

func TestGetUsers(t *testing.T) {
	r, _ := setupUserTest()

	// テストユーザーを作成
	authService := services.NewAuthService(testutils.TestDB)
	testUsers := []models.User{
		{Email: "user1@example.com", Password: "password123"},
		{Email: "user2@example.com", Password: "password123"},
	}

	var tokens []string
	for _, user := range testUsers {
		err := authService.SignUp(&user)
		assert.NoError(t, err)
		token, err := authService.Login(user.Email, "password123")
		assert.NoError(t, err)
		tokens = append(tokens, token.Token)
	}

	tests := []struct {
		name           string
		token          string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "Successful Get Users",
			token:          tokens[0],
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "Unauthorized Access",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
			expectedCount:  0,
		},
		{
			name:           "Invalid Token",
			token:          "invalid_token",
			expectedStatus: http.StatusUnauthorized,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/users", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response []models.UserResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response, tt.expectedCount)

				// レスポンスの内容を確認
				for _, user := range response {
					assert.NotEmpty(t, user.ID)
					assert.NotEmpty(t, user.Email)
				}
			}
		})
	}
}
