package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-gin-gorm-minimum/handlers"
	"go-gin-gorm-minimum/models"
	"go-gin-gorm-minimum/services"
	"go-gin-gorm-minimum/testutils"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gin-gorm-minimum/middlewares"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	db               *gorm.DB
	micropostHandler *handlers.MicropostHandler
	authService      *services.AuthService
)

func TestCreateMicropost(t *testing.T) {
	r, handler, _ := testutils.SetupMicropostHandler()

	// ルートの設定
	r.POST("/microposts", middlewares.AuthMiddleware(), handler.CreateMicropost)

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
				_, token, err = testutils.CreateTestUser()
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
