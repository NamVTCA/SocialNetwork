package notification


import (
	"net/http"


	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationHandler struct {
	service NotificationService
}

func NewNotificationHandler(service NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (h *NotificationHandler) GetUserNotifications(c *gin.Context) {
	userIDStr := c.GetString("userID") // bạn lấy từ middleware JWT
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	notifications, err := h.service.GetUserNotifications(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get notifications"})
		return
	}
	c.JSON(http.StatusOK, notifications)
}

func (h *NotificationHandler) ReadNotification(c *gin.Context) {
	notifID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(notifID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}
	if err := h.service.ReadNotification(c.Request.Context(), objID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark as read"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Marked as read"})
}