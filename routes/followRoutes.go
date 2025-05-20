package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"socialnetwork/internal/follow"
)

// RegisterFollowRoutes mounts follow routes.
func RegisterFollowRoutes(rg *gin.RouterGroup, db *mongo.Database) {
	h := follow.NewFollowHandler(follow.NewFollowService(follow.NewFollowRepository(db)))
	f := rg.Group("/follows")
	f.POST("/", h.Follow)
	f.DELETE("/:id", h.Unfollow)
	f.GET("/:id/followers", h.ListFollowers)
	f.GET("/:id/following", h.ListFollowing)
}