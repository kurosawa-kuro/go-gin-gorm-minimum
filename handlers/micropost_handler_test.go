package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-gin-gorm-minimum/handlers"
	"go-gin-gorm-minimum/models"
	"go-gin-gorm-minimum/services"
	"go-gin-gorm-minimum/testutils"
	"go-gin-gorm-minimum/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gin-gorm-minimum/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	db               *gorm.DB
	micropostHandler *handlers.MicropostHandler
	authService      *services.AuthService
)

func setupMicropostTest() (*gin.Engine, *handlers.MicropostHandler) {
	if err := testutils.CleanupDatabase(testutils.TestDB); err != nil {
		panic(err)
	}

	micropostService := services.NewMicropostService(testutils.TestDB)
	micropostHandler = handlers.NewMicropostHandler(micropostService)
	authService = services.NewAuthService(testutils.TestDB)

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// JWT middleware setup

	return r, micropostHandler
}

func createTestUser() (*models.User, string, error) {
	user := &models.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	if err := authService.SignUp(user); err != nil {
		return nil, "", err
	}

	token, err := utils.GenerateJWTToken(*user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func TestCreateMicropost(t *testing.T) {
	r, handler := setupMicropostTest()
	authMiddleware := middlewares.AuthMiddleware()
	r.POST("/microposts", authMiddleware, handler.CreateMicropost)

	tests := []struct {
		name           string
		setupAuth      bool
		request        models.MicropostRequest
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:      "Successful Creation",
			setupAuth: true,
			request: models.MicropostRequest{
				Title: "Test Micropost",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.MicropostResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Test Micropost", response.Title)
				assert.NotZero(t, response.ID)
				assert.NotEmpty(t, response.CreatedAt)
			},
		},
		{
			name:      "Unauthorized Access",
			setupAuth: false,
			request: models.MicropostRequest{
				Title: "Test Micropost",
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
			name:      "Empty Title",
			setupAuth: true,
			request: models.MicropostRequest{
				Title: "",
			},
			expectedStatus: http.StatusBadRequest,
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
			err := testutils.CleanupDatabase(testutils.TestDB)
			assert.NoError(t, err)

			var token string
			if tt.setupAuth {
				_, token, err = createTestUser()
				assert.NoError(t, err)
			}

			// リクエストの作成
			requestBody, err := json.Marshal(tt.request)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/microposts", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			if tt.setupAuth {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}
