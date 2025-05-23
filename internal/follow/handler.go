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
	userIDStr := c.MustGet("userID").(string)
	followerID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID in token"})
		return
	}

	followingIDHex := c.Param("id")
	followingID, err := primitive.ObjectIDFromHex(followingIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target user ID"})
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
	userIDStr := c.MustGet("userID").(string)
	followerID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID in token"})
		return
	}

	followingIDHex := c.Param("id")
	followingID, err := primitive.ObjectIDFromHex(followingIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target user ID"})
		return
	}

	// 🚫 Không cho unfollow chính mình
	if followerID == followingID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You cannot unfollow yourself"})
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
	userIDStr := c.MustGet("userID").(string)

	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	// Nếu người dùng đang xem chính mình
	if id.Hex() == userIDStr {
		followers, err := h.followService.GetFollowers(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"count":     len(followers),
			"followers": followers,
		})
		return
	}

	// Nếu không phải chính mình -> chỉ trả về số lượng
	count, err := h.followService.GetFollowerCount(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"count": count,
	})
}



// ListFollowing handles GET /:id/following requests
func (h *FollowHandler) ListFollowing(c *gin.Context) {
	idHex := c.Param("id")
	userIDStr := c.MustGet("userID").(string)

	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	if id.Hex() == userIDStr {
		following, err := h.followService.GetFollowing(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"count":     len(following),
			"following": following,
		})
		return
	}

	// Người khác -> chỉ trả về số lượng
	count, err := h.followService.GetFollowingCount(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"count": count,
	})
}
