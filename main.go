package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	_ "go-gin-gorm-minimum/docs"

	"go-gin-gorm-minimum/middlewares"
	"go-gin-gorm-minimum/models"
	"go-gin-gorm-minimum/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// @title           API
// @version         1.0
// @description     This is a sample server.
// @host           localhost:8080
// @BasePath       /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func init() {
	// データベース接続
	var err error
	dsn := "host=localhost user=postgres password=postgres dbname=web_app_db_integration_go port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// マイグレーション
	db.AutoMigrate(&models.Micropost{}, &models.User{})
}

// CreateMicropost godoc
// @Summary      Create new micropost
// @Description  Create a new micropost with the given title
// @Tags         microposts
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        micropost body models.Micropost true "Micropost object"
// @Success      201  {object}  models.Micropost
// @Router       /api/v1/microposts [post]
func CreateMicropost(c *gin.Context) {
	userID, _ := c.Get("user_id")
	fmt.Println("userID:", userID)
	var micropost models.Micropost
	if err := c.ShouldBindJSON(&micropost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Create(&micropost)
	c.JSON(http.StatusCreated, micropost)
}

// GetMicroposts godoc
// @Summary      List microposts
// @Description  get all microposts
// @Tags         microposts
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {array}   models.Micropost
// @Router       /api/v1/microposts [get]
func GetMicroposts(c *gin.Context) {
	var microposts []models.Micropost
	db.Find(&microposts)
	c.JSON(http.StatusOK, microposts)
}

// GetMicropost godoc
// @Summary      Get micropost by ID
// @Description  get micropost by ID
// @Tags         microposts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Micropost ID"
// @Success      200  {object}  models.Micropost
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/microposts/{id} [get]
func GetMicropost(c *gin.Context) {
	var micropost models.Micropost
	if err := db.First(&micropost, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, micropost)
}

// SignupUser godoc
// @Summary      Signup user
// @Description  Signup user with the given information
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user body models.User true "User object" default({"email":"user1@example.com","password":"password123","role":"user","avatar_path":"/avatars/default.png"})
// @Success      201  {object}  models.UserResponse
// @Router       /api/v1/auth/signup [post]
func SignupUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// メールアドレスの重複チェック
	var existingUser models.User
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = hashedPassword

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user.ToResponse())
}

// LoginUser godoc
// @Summary      Login user
// @Description  Login user with the given email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user body models.User true "User object" default({"email":"user1@example.com","password":"password123"})
// @Success      200  {object}  models.LoginResponse
// @Router       /api/v1/auth/login [post]
func LoginUser(c *gin.Context) {
	var loginUser models.User
	var storedUser models.User

	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Where("email = ?", loginUser.Email).First(&storedUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := utils.CheckPassword(storedUser.Password, loginUser.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   storedUser.ID,
		"email": storedUser.Email,
		"exp":   time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	fmt.Println("tokenString:Bearer", tokenString)

	response := models.LoginResponse{
		Token:        tokenString,
		UserResponse: storedUser.ToResponse(),
	}

	c.JSON(http.StatusOK, response)
}

// GetUsers godoc
// @Summary      List users
// @Description  get all users
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {array}   models.UserResponse
// @Router       /api/v1/users [get]
func GetUsers(c *gin.Context) {
	var users []models.User
	db.Find(&users)

	var response []models.UserResponse
	for _, user := range users {
		response = append(response, user.ToResponse())
	}

	c.JSON(http.StatusOK, response)
}

// GetUser godoc
// @Summary      Get user by ID
// @Description  get user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  models.UserResponse
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/users/{id} [get]
func GetUser(c *gin.Context) {
	var user models.User
	if err := db.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// GetMe godoc
// @Summary      Get current user
// @Description  get current user information from token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {object}  models.UserResponse
// @Failure      401  {object}  map[string]string
// @Router       /api/v1/auth/me [get]
func GetMe(c *gin.Context) {
	userID, _ := c.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

func main() {
	r := gin.Default()

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	v1 := r.Group("/api/v1")
	{
		microposts := v1.Group("/microposts")
		microposts.Use(middlewares.AuthMiddleware())
		{
			microposts.POST("", CreateMicropost)
			microposts.GET("", GetMicroposts)
			microposts.GET("/:id", GetMicropost)
		}

		users := v1.Group("/users")
		users.Use(middlewares.AuthMiddleware())
		{
			users.GET("", GetUsers)
			users.GET("/:id", GetUser)
		}

		auth := v1.Group("/auth")
		{
			auth.POST("/signup", SignupUser)
			auth.POST("/login", LoginUser)
			auth.GET("/me", middlewares.AuthMiddleware(), GetMe)
		}
	}

	r.Run(":8080")
}
