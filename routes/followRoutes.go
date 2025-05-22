package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"socialnetwork/internal/follow"
	"socialnetwork/internal/notification"
	"socialnetwork/internal/user"
	"socialnetwork/pkg/middleware"
)

func FollowRoutes(rg *gin.RouterGroup, db *mongo.Database, followHandler *follow.FollowHandler) {
	followRepo := follow.NewFollowRepository(db)
	userRepo := user.NewRepository(db)
	notifRepo := notification.NewNotificationRepository(db)

	service := follow.NewFollowService(followRepo, userRepo, notifRepo)
	handler := follow.NewFollowHandler(service)

	followRoutes := rg.Group("/follows", middleware.JWTAuthMiddleware())
	{
		followRoutes.POST("/:id", handler.FollowUser)              // follow user with id
		followRoutes.DELETE("/:id", handler.UnfollowUser)          // unfollow user with id
		followRoutes.GET("/:id/followers", handler.ListFollowers)
		followRoutes.GET("/:id/following", handler.ListFollowing)
	}
}
