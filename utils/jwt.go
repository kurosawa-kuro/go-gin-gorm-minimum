package utils

import (
	"fmt"
	"go-gin-gorm-minimum/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// GenerateJWTToken creates a new JWT token for the given user
func GenerateJWTToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour).Unix(),
	})
	signedToken, err := token.SignedString(jwtSecret)
	fmt.Println("token:Bearer", signedToken)
	return signedToken, err
}

// ParseJWTToken validates and parses the given JWT token
func ParseJWTToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
}
