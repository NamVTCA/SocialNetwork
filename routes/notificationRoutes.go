package routes

import (
	"socialnetwork/internal/notification"

	"github.com/gin-gonic/gin"
)

func RegisterNotificationRoutes(rg *gin.RouterGroup, h *notification.NotificationHandler) {
	notif := rg.Group("/notifications")
	notif.GET("/", h.GetUserNotifications)
	notif.PUT("/:id/read", h.ReadNotification)
}
