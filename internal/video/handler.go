package video

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"socialnetwork/models"
)

type VideoHandler struct {
	videoService VideoService
}

func NewVideoHandler(videoService VideoService) *VideoHandler {
	return &VideoHandler{videoService}
}

func (h *VideoHandler) CreateVideo(c *gin.Context) {
	var video models.Video
	if err := c.ShouldBindJSON(&video); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr := c.MustGet("userID").(string)
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}
	video.OwnerID = userID

	// Validate visibility (nếu không có thì mặc định public)
	if video.Visibility != "private" && video.Visibility != "public" {
		video.Visibility = "public"
	}

	if err := h.videoService.CreateVideo(c.Request.Context(), &video); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, video)
}



func (h *VideoHandler) GetVideoByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video ID"})
		return
	}

	video, err := h.videoService.GetVideoByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "video not found"})
		return
	}

	// Lấy userID (nếu đã đăng nhập)
	var userID primitive.ObjectID
	userIDStr, exists := c.Get("userID")
	if exists {
		userID, _ = primitive.ObjectIDFromHex(userIDStr.(string))
	}

	// Nếu video private mà không phải chủ sở hữu thì trả lỗi
	if video.Visibility == "private" && video.OwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bạn không có quyền xem video này"})
		return
	}

	c.JSON(http.StatusOK, video)
}


func (h *VideoHandler) GetVideosByOwner(c *gin.Context) {
	ownerIDStr := c.Param("ownerID") // giả sử url có tham số ownerID
	ownerID, err := primitive.ObjectIDFromHex(ownerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid owner ID"})
		return
	}

	requesterIDStr, exists := c.Get("userID")
	var requesterID primitive.ObjectID
	if exists {
		requesterID, _ = primitive.ObjectIDFromHex(requesterIDStr.(string))
	}

	// Nếu requester là owner thì lấy tất cả video (cả public và private)
	// Nếu không thì chỉ lấy video public
	var videos []models.Video
	if requesterID == ownerID {
		videos, err = h.videoService.GetVideosByOwner(c.Request.Context(), ownerID)
	} else {
		videos, err = h.videoService.GetPublicVideosByOwner(c.Request.Context(), ownerID)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch videos"})
		return
	}

	c.JSON(http.StatusOK, videos)
}


func (h *VideoHandler) DeleteVideo(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video ID"})
		return
	}

	ownerID := c.MustGet("userID").(primitive.ObjectID)
	if err := h.videoService.DeleteVideo(c.Request.Context(), id, ownerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete video"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "video deleted"})
}
