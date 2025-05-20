package follow


import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	
)

// FollowHandler provides follow endpoints.
type FollowHandler struct {
	service FollowService
}

func NewFollowHandler(s FollowService) *FollowHandler {
	return &FollowHandler{service: s}
}

func (h *FollowHandler) Follow(c *gin.Context) {
	uid := c.GetString("userID")
	follower, _ := primitive.ObjectIDFromHex(uid)
	var body struct { Following string `json:"following" binding:"required"` }
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	followee, _ := primitive.ObjectIDFromHex(body.Following)
	if err := h.service.Follow(c.Request.Context(), follower, followee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "followed"})
}

func (h *FollowHandler) Unfollow(c *gin.Context) {
	uid := c.GetString("userID")
	follower, _ := primitive.ObjectIDFromHex(uid)
	followeeID := c.Param("id")
	followee, _ := primitive.ObjectIDFromHex(followeeID)
	if err := h.service.Unfollow(c.Request.Context(), follower, followee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "unfollowed"})
}

func (h *FollowHandler) ListFollowers(c *gin.Context) {
	userID := c.Param("id")
	uid, _ := primitive.ObjectIDFromHex(userID)
	list, _ := h.service.GetFollowers(c.Request.Context(), uid)
	c.JSON(http.StatusOK, list)
}

func (h *FollowHandler) ListFollowing(c *gin.Context) {
	userID := c.Param("id")
	uid, _ := primitive.ObjectIDFromHex(userID)
	list, _ := h.service.GetFollowing(c.Request.Context(), uid)
	c.JSON(http.StatusOK, list)
}