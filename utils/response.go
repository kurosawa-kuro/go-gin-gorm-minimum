package utils

import "github.com/gin-gonic/gin"

// Common error responses
var (
	ErrInvalidCredentials = gin.H{"error": "Invalid email or password"}
	ErrUnauthorized       = gin.H{"error": "Unauthorized"}
	ErrRecordNotFound     = gin.H{"error": "Record not found!"}
)
