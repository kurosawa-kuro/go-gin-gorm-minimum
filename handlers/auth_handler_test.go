package handlers_test

import (
	"go-gin-gorm-minimum/testutils"
	"testing"
)

func TestLoginUser(t *testing.T) {
	r, authHandler := testutils.SetupAuthHandler()

	// ルートの設定
	r.POST("/auth/login", authHandler.LoginUser)

	// テストの残りの部分...
}
