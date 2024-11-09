package main

import (
	"net/http"
	"os"
	"time"

	_ "go-gin-gorm-minimum/docs"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// User モデル定義
type User struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Email      string    `json:"email" gorm:"uniqueIndex;not null" binding:"required,email" example:"user@example.com"`
	Password   string    `json:"password" gorm:"not null" binding:"required,min=6" example:"password123"`
	Role       string    `json:"role" gorm:"default:'user'" example:"user"`
	AvatarPath string    `json:"avatar_path" example:"/avatars/default.png"`
	CreatedAt  time.Time `json:"-" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"-" gorm:"autoUpdateTime"`
}

// Micropost モデル定義
type Micropost struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" binding:"required" example:"マイクロポストのタイトル"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// @title           API
// @version         1.0
// @description     This is a sample server.
// @host           localhost:8080
// @BasePath       /

func init() {
	// データベース接続
	var err error
	dsn := "host=localhost user=postgres password=postgres dbname=web_app_db_integration_go port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// マイグレーション
	db.AutoMigrate(&Micropost{}, &User{}) // User モデルを追加
}

// CreateMicropost godoc
// @Summary      Create new micropost
// @Description  Create a new micropost with the given title
// @Tags         microposts
// @Accept       json
// @Produce      json
// @Param        micropost body Micropost true "Micropost object"
// @Success      201  {object}  Micropost
// @Router       /api/v1/microposts [post]
func CreateMicropost(c *gin.Context) {
	var micropost Micropost
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
// @Success      200  {array}   Micropost
// @Router       /api/v1/microposts [get]
func GetMicroposts(c *gin.Context) {
	var microposts []Micropost
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
// @Success      200  {object}  Micropost
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/microposts/{id} [get]
func GetMicropost(c *gin.Context) {
	var micropost Micropost
	if err := db.First(&micropost, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, micropost)
}

// CreateUser godoc
// @Summary      Create new user
// @Description  Create a new user with the given information
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user body User true "User object"
// @Success      201  {object}  User
// @Router       /api/v1/users [post]
func CreateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)
	db.Create(&user)

	// レスポンス用の構造体を作成
	response := struct {
		ID         uint      `json:"id"`
		Email      string    `json:"email"`
		Role       string    `json:"role"`
		AvatarPath string    `json:"avatar_path"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	}{
		ID:         user.ID,
		Email:      user.Email,
		Role:       user.Role,
		AvatarPath: user.AvatarPath,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// Login User
// @Summary      Login user
// @Description  Login user with the given email and password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user body User true "User object"
// @Success      200  {object}  User
// @Router       /api/v1/users/login [post]
func LoginUser(c *gin.Context) {
	var loginUser User
	var storedUser User

	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// データベースからユーザーを検索
	if err := db.Where("email = ?", loginUser.Email).First(&storedUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// パスワードの比較
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(loginUser.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// token関連
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   storedUser.ID,
		"email": storedUser.Email,
		"exp":   time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// パスワードを除外したレスポンスを作成
	response := struct {
		Token      string    `json:"token"`
		ID         uint      `json:"id"`
		Email      string    `json:"email"`
		Role       string    `json:"role"`
		AvatarPath string    `json:"avatar_path"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	}{
		Token:      tokenString,
		ID:         storedUser.ID,
		Email:      storedUser.Email,
		Role:       storedUser.Role,
		AvatarPath: storedUser.AvatarPath,
		CreatedAt:  storedUser.CreatedAt,
		UpdatedAt:  storedUser.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// GetUsers godoc
// @Summary      List users
// @Description  get all users
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {array}   User
// @Router       /api/v1/users [get]
func GetUsers(c *gin.Context) {
	var users []User
	db.Find(&users)
	c.JSON(http.StatusOK, users)
}

// GetUser godoc
// @Summary      Get user by ID
// @Description  get user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  User
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/users/{id} [get]
func GetUser(c *gin.Context) {
	var user User
	if err := db.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func main() {
	r := gin.Default()

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	v1 := r.Group("/api/v1")
	{
		microposts := v1.Group("/microposts")
		{
			microposts.POST("", CreateMicropost)
			microposts.GET("", GetMicroposts)
			microposts.GET("/:id", GetMicropost)
		}

		users := v1.Group("/users")
		{
			users.POST("", CreateUser)
			users.POST("/login", LoginUser)
			users.GET("", GetUsers)
			users.GET("/:id", GetUser)
		}
	}

	r.Run(":8080")
}
