package main

import (
	"net/http"

	_ "go-gin-gorm-minimum/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// Micropost モデル定義
type Micropost struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Title string `json:"title" binding:"required" example:"マイクロポストのタイトル"`
}

// @title           Micropost API
// @version         1.0
// @description     This is a micropost server.
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
	db.AutoMigrate(&Micropost{})
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
	}

	r.Run(":8080")
}
