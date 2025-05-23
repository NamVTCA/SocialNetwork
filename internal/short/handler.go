package short

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"socialnetwork/models"
)

type ShortHandler struct {
	service ShortService
}

func NewShortHandler(service ShortService) *ShortHandler {
	return &ShortHandler{service}
}

func (h *ShortHandler) CreateShort(c *gin.Context) {
	var sh models.Short
	if err := c.ShouldBindJSON(&sh); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sh.OwnerID = c.MustGet("userID").(primitive.ObjectID)
	if err := h.service.CreateShort(c.Request.Context(), &sh); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sh)
}

func (h *ShortHandler) GetShortByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid short ID"})
		return
	}

	sh, err := h.service.GetShortByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "short not found"})
		return
	}

	c.JSON(http.StatusOK, sh)
}

func (h *ShortHandler) GetShortsByOwner(c *gin.Context) {
	ownerID := c.MustGet("userID").(primitive.ObjectID)
	list, err := h.service.GetShortsByOwner(c.Request.Context(), ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch shorts"})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *ShortHandler) DeleteShort(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid short ID"})
		return
	}

	ownerID := c.MustGet("userID").(primitive.ObjectID)
	if err := h.service.DeleteShort(c.Request.Context(), id, ownerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete short"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "short deleted"})
}

func (h *ShortHandler) GetPublicShortsByOwner(c *gin.Context) {
	ownerIDStr := c.Param("ownerID")
	ownerID, err := primitive.ObjectIDFromHex(ownerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ownerID"})
		return
	}

	shorts, err := h.service.GetPublicShortsByOwner(c.Request.Context(), ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get public shorts"})
		return
	}

	c.JSON(http.StatusOK, shorts)
}
