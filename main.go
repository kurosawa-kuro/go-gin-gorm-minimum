package main

import (
	"net/http"

	_ "go-gin-gorm-minimum/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Gin Swagger Example API
// @version         1.0
// @description     This is a sample server.
// @host           localhost:8080
// @BasePath       /

// HelloWorld godoc
// @Summary      Hello world endpoint
// @Description  get hello world message
// @Tags         example
// @Accept       json
// @Produce      json
// @Success      200  {object}  HelloResponse
// @Router       / [get]
func HelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

// HelloResponse represents the response structure
type HelloResponse struct {
	Message string `json:"message" example:"Hello, World!"`
}

func main() {
	r := gin.Default()

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routes
	r.GET("/", HelloWorld)

	r.Run(":8080")
}
