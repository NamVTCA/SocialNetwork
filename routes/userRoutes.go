package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"socialnetwork/pkg/middleware"
	"socialnetwork/internal/user"
)

func UserRoutes(r *gin.Engine, db *mongo.Database, handler *user.Handler) {
	userRoutes := r.Group("/users")
	{
		userRoutes.GET("/", handler.GetUsers)
		userRoutes.GET("/me", middleware.JWTAuthMiddleware(), handler.GetMe)
	}
}
