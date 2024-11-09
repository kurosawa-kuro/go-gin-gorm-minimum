package utils

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

// 許可する画像の拡張子
var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}

// IsValidImageFile checks if the uploaded file is a valid image
func IsValidImageFile(file *multipart.FileHeader) bool {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	return allowedExtensions[ext]
}

// GenerateAvatarFilename generates a unique filename for the avatar
func GenerateAvatarFilename(userID uint, originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	timestamp := time.Now().Unix()
	return fmt.Sprintf("user_%d_%d%s", userID, timestamp, ext)
}
