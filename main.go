package main

import (
	"net/http"
	"time"

	_ "go-gin-gorm-minimum/docs"

	"github.com/gin-gonic/gin"
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
	c.JSON(http.StatusCreated, user)
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
			users.GET("", GetUsers)
			users.GET("/:id", GetUser)
		}
	}

	r.Run(":8080")
}
