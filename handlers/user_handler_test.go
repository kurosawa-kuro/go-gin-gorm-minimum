package handlers_test

import (
	"bytes"
	"encoding/json"
	"go-gin-gorm-minimum/handlers"
	"go-gin-gorm-minimum/infra"
	"go-gin-gorm-minimum/middlewares"
	"go-gin-gorm-minimum/models"
	"go-gin-gorm-minimum/services"
	"go-gin-gorm-minimum/testutils"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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

func TestUpdateAvatar(t *testing.T) {
	r, userHandler := setupUserTest()

	// Add route for avatar update
	auth := r.Group("/")
	auth.Use(middlewares.AuthMiddleware())
	auth.PUT("/users/avatar", userHandler.UpdateAvatar)

	// Create test user and get token
	authService := services.NewAuthService(testutils.TestDB)
	testUser := models.User{
		Email:    "avatar.test@example.com",
		Password: "password123",
	}
	err := authService.SignUp(&testUser)
	assert.NoError(t, err)
	token, err := authService.Login(testUser.Email, "password123")
	assert.NoError(t, err)

	// Create test image file
	testImagePath := filepath.Join("testdata", "test.jpg")
	err = os.MkdirAll("testdata", 0755)
	assert.NoError(t, err)
	testImage, err := os.Create(testImagePath)
	assert.NoError(t, err)
	defer os.Remove(testImagePath)
	defer testImage.Close()

	// Write some dummy image data
	_, err = testImage.Write([]byte("fake image content"))
	assert.NoError(t, err)
	testImage.Close()

	tests := []struct {
		name           string
		token          string
		setupFile      func() (*bytes.Buffer, *multipart.Writer)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:  "Successful Avatar Update",
			token: token.Token,
			setupFile: func() (*bytes.Buffer, *multipart.Writer) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, err := writer.CreateFormFile("avatar", "test.jpg")
				assert.NoError(t, err)

				file, err := os.Open(testImagePath)
				assert.NoError(t, err)
				defer file.Close()

				_, err = io.Copy(part, file)
				assert.NoError(t, err)
				writer.Close()
				return body, writer
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.UserResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response.AvatarPath, "/uploads/avatars/")
				assert.Contains(t, response.AvatarPath, ".jpg")
			},
		},
		{
			name:  "Unauthorized Access",
			token: "",
			setupFile: func() (*bytes.Buffer, *multipart.Writer) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				writer.Close()
				return body, writer
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:  "No File Uploaded",
			token: token.Token,
			setupFile: func() (*bytes.Buffer, *multipart.Writer) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				writer.Close()
				return body, writer
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, writer := tt.setupFile()

			req := httptest.NewRequest(http.MethodPut, "/users/avatar", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}

	// Cleanup test uploads directory after all tests
	err = os.RemoveAll("uploads")
	assert.NoError(t, err)
}
