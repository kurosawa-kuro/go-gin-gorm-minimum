// @title           Go Gin GORM Minimum API
// @version         1.0
// @description     A minimal Go REST API with Gin and GORM.
// @host           localhost:8080
// @BasePath       /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer {token}

package main

import (
	_ "go-gin-gorm-minimum/docs"
	"go-gin-gorm-minimum/handlers"
	"go-gin-gorm-minimum/infra"
	"go-gin-gorm-minimum/middlewares"
	"go-gin-gorm-minimum/models"
	"go-gin-gorm-minimum/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// Initialize database connection and migrations
func init() {
	var err error
	db, err = gorm.Open(postgres.Open("postgresql://postgres:postgres@localhost:5432/web_app_db_integration_go?sslmode=disable"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Micropost{}, &models.User{})
}

// Router setup
type Router struct {
	auth      *handlers.AuthHandler
	user      *handlers.UserHandler
	micropost *handlers.MicropostHandler
}

func NewRouter(db *gorm.DB) *Router {
	userService := services.NewUserService(db)
	authService := services.NewAuthService(db)
	micropostService := services.NewMicropostService(db)

	return &Router{
		auth:      handlers.NewAuthHandler(authService, userService),
		user:      handlers.NewUserHandler(userService),
		micropost: handlers.NewMicropostHandler(micropostService),
	}
}

func (router *Router) Setup(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		router.setupMicropostRoutes(v1.Group("/microposts"))
		router.setupUserRoutes(v1.Group("/users"))
		router.setupAuthRoutes(v1.Group("/auth"))
	}
}

func (router *Router) setupMicropostRoutes(group *gin.RouterGroup) {
	group.Use(middlewares.AuthMiddleware())
	group.POST("", router.micropost.CreateMicropost)
	group.GET("", router.micropost.GetMicroposts)
	group.GET("/:id", router.micropost.GetMicropost)
}

func (router *Router) setupUserRoutes(group *gin.RouterGroup) {
	group.Use(middlewares.AuthMiddleware())
	group.GET("", router.user.GetUsers)
	group.GET("/:id", router.user.GetUser)
	group.PUT("/avatar", router.user.UpdateAvatar)
}

func (router *Router) setupAuthRoutes(group *gin.RouterGroup) {
	group.POST("/signup", router.auth.SignupUser)
	group.POST("/login", router.auth.LoginUser)
	group.GET("/me", middlewares.AuthMiddleware(), router.auth.GetMe)
}

func main() {
	db := infra.SetupDB()
	router := NewRouter(db)

	r := gin.Default()
	router.Setup(r)
	r.Run(":8080")
}
