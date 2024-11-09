// @title           Go Gin GORM Minimum API
// @version         1.0
// @description     A minimal Go REST API with Gin and GORM.
// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer {token}

package main

import (
	"net/http"

	_ "go-gin-gorm-minimum/docs"

	"go-gin-gorm-minimum/middlewares"
	"go-gin-gorm-minimum/models"
	"go-gin-gorm-minimum/utils"

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
	db, err = gorm.Open(postgres.Open(utils.DSN), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Micropost{}, &models.User{})
}

// SignupUser godoc
// @Summary      Signup user
// @Description  Signup user with the given information
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user body models.User true "User object" default({"email":"user1@example.com","password":"password123","role":"user","avatar_path":"/avatars/default.png"})
// @Success      201  {object}  models.UserResponse
// @Router       /auth/signup [post]
func SignupUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := utils.FindUserByEmail(db, user.Email); err == nil {
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
// @Param        user body models.LoginRequest true "Login credentials"
// @Success      200  {object}  models.LoginResponse
// @Router       /auth/login [post]
// @Security BearerAuth
func LoginUser(c *gin.Context) {
	var loginReq models.LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	storedUser, err := utils.FindUserByEmail(db, loginReq.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrInvalidCredentials)
		return
	}

	if err := utils.CheckPassword(storedUser.Password, loginReq.Password); err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrInvalidCredentials)
		return
	}

	tokenString, err := utils.GenerateJWTToken(storedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token:        tokenString,
		UserResponse: storedUser.ToResponse(),
	})
}

// GetMe godoc
// @Summary      Get current user
// @Description  get current user information from token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.UserResponse
// @Failure      401  {object}  map[string]string
// @Router       /auth/me [get]
// @Security BearerAuth
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

// GetUsers godoc
// @Summary      List users
// @Description  get all users
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.UserResponse
// @Router       /users [get]
// @Security BearerAuth
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
// @Security     BearerAuth
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  models.UserResponse
// @Failure      404  {object}  map[string]string
// @Router       /users/{id} [get]
// @Security BearerAuth
func GetUser(c *gin.Context) {
	var user models.User
	if err := db.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// CreateMicropost godoc
// @Summary      Create new micropost
// @Description  Create a new micropost with the given title
// @Tags         microposts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        micropost body models.Micropost true "Micropost object"
// @Success      201  {object}  models.Micropost
// @Router       /microposts [post]
// @Security BearerAuth
func CreateMicropost(c *gin.Context) {
	userID, _ := c.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var micropost models.Micropost
	if err := c.ShouldBindJSON(&micropost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the UserID from the authenticated user
	micropost.UserID = userID.(uint)

	if err := db.Create(&micropost).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create micropost"})
		return
	}
	c.JSON(http.StatusCreated, micropost)
}

// GetMicroposts godoc
// @Summary      List microposts
// @Description  get all microposts
// @Tags         microposts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.Micropost
// @Router       /microposts [get]
// @Security BearerAuth
func GetMicroposts(c *gin.Context) {
	var microposts []models.Micropost
	db.Preload("User").Find(&microposts)
	c.JSON(http.StatusOK, microposts)
}

// GetMicropost godoc
// @Summary      Get micropost by ID
// @Description  get micropost by ID
// @Tags         microposts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Micropost ID"
// @Success      200  {object}  models.Micropost
// @Failure      404  {object}  map[string]string
// @Router       /microposts/{id} [get]
// @Security BearerAuth
func GetMicropost(c *gin.Context) {
	var micropost models.Micropost
	if err := db.First(&micropost, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, micropost)
}

func main() {
	r := gin.Default()
	setupRoutes(r)
	r.Run(":8080")
}

// setupRoutes configures all the routes for the application
func setupRoutes(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		setupMicropostRoutes(v1.Group("/microposts"))
		setupUserRoutes(v1.Group("/users"))
		setupAuthRoutes(v1.Group("/auth"))
	}
}

func setupMicropostRoutes(group *gin.RouterGroup) {
	group.Use(middlewares.AuthMiddleware())
	group.POST("", CreateMicropost)
	group.GET("", GetMicroposts)
	group.GET("/:id", GetMicropost)
}

func setupUserRoutes(group *gin.RouterGroup) {
	group.Use(middlewares.AuthMiddleware())
	group.GET("", GetUsers)
	group.GET("/:id", GetUser)
}

func setupAuthRoutes(group *gin.RouterGroup) {
	group.POST("/signup", SignupUser)
	group.POST("/login", LoginUser)
	group.GET("/me", middlewares.AuthMiddleware(), GetMe)
}
