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

	video.OwnerID = c.MustGet("userID").(primitive.ObjectID)
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

	c.JSON(http.StatusOK, video)
}

func (h *VideoHandler) GetVideosByOwner(c *gin.Context) {
	ownerID := c.MustGet("userID").(primitive.ObjectID)
	videos, err := h.videoService.GetVideosByOwner(c.Request.Context(), ownerID)
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
