package follow

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FollowHandler struct {
	followService FollowService
}

func NewFollowHandler(followService FollowService) *FollowHandler {
	return &FollowHandler{followService: followService}
}

func (h *FollowHandler) FollowUser(c *gin.Context) {
	followerID := c.MustGet("userID").(primitive.ObjectID) // giả sử userID được lấy từ middleware auth

	followingIDHex := c.Param("id")
	followingID, err := primitive.ObjectIDFromHex(followingIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	err = h.followService.FollowUser(c.Request.Context(), followerID, followingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "followed"})
}

func (h *FollowHandler) UnfollowUser(c *gin.Context) {
	followerID := c.MustGet("userID").(primitive.ObjectID)

	followingIDHex := c.Param("id")
	followingID, err := primitive.ObjectIDFromHex(followingIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	err = h.followService.UnfollowUser(c.Request.Context(), followerID, followingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "unfollowed"})
}

// ListFollowers handles GET /:id/followers requests
func (h *FollowHandler) ListFollowers(c *gin.Context) {
	idHex := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	followers, err := h.followService.GetFollowers(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, followers)
}


// ListFollowing handles GET /:id/following requests
func (h *FollowHandler) ListFollowing(c *gin.Context) {
	idHex := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	followings, err := h.followService.GetFollowing(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, followings)
}