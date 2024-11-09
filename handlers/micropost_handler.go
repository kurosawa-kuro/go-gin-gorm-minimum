package handlers

import (
	"net/http"

	"go-gin-gorm-minimum/models"
	"go-gin-gorm-minimum/services"

	"github.com/gin-gonic/gin"
)

type MicropostHandler struct {
	micropostService *services.MicropostService
}

func NewMicropostHandler(micropostService *services.MicropostService) *MicropostHandler {
	return &MicropostHandler{micropostService: micropostService}
}

// CreateMicropost godoc
// @Summary      Create new micropost
// @Description  Create a new micropost with the given title
// @Tags         microposts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        micropost body models.MicropostRequest true "Micropost object"
// @Success      201  {object}  models.MicropostResponse
// @Router       /microposts [post]
func (h *MicropostHandler) CreateMicropost(c *gin.Context) {
	userID, _ := c.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.MicropostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	micropost := models.Micropost{
		Title:  req.Title,
		UserID: userID.(uint),
	}

	if err := h.micropostService.Create(&micropost); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create micropost"})
		return
	}

	c.JSON(http.StatusCreated, micropost.ToResponse())
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
func (h *MicropostHandler) GetMicroposts(c *gin.Context) {
	microposts, err := h.micropostService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch microposts"})
		return
	}
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
func (h *MicropostHandler) GetMicropost(c *gin.Context) {
	micropost, err := h.micropostService.GetByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, micropost)
}
