package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"socialnetwork/internal/user"
	"socialnetwork/pkg/middleware"
)

func UserRoutes(r *gin.Engine, db *mongo.Database, handler *user.Handler) {
	userRoutes := r.Group("/users")
	{
		userRoutes.GET("/", handler.GetUsers)
		userRoutes.GET("/me", middleware.JWTAuthMiddleware(), handler.GetMe)
		userRoutes.PUT("/me", middleware.JWTAuthMiddleware(), handler.UpdateMe)
		userRoutes.POST("/change-password", middleware.JWTAuthMiddleware(), handler.ChangePassword)
		userRoutes.POST("/forgot-password", handler.ForgotPassword)
		userRoutes.POST("/reset-password", handler.ResetPassword)
		userRoutes.POST("/change-email", middleware.JWTAuthMiddleware(), handler.ChangeEmailRequest)
		userRoutes.POST("/verify-email", middleware.JWTAuthMiddleware(), handler.VerifyEmailRequest)
	}
}
