package routes

import (
    "socialnetwork/internal/comment"
    "socialnetwork/pkg/middleware"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/mongo"
)

func CommentRoutes(r *gin.Engine, db *mongo.Database, handler *comment.CommentHandler) {
    commentGroup := r.Group("/comments")


    commentGroup.POST("/:postID", handler.CreateComment)
    commentGroup.GET("/:postID", handler.GetCommentsByPost)
    commentGroup.PUT("/:id", middleware.JWTAuthMiddleware(), handler.UpdateComment)
    commentGroup.DELETE("/:id", middleware.JWTAuthMiddleware(), handler.DeleteComment)
	commentGroup.PUT("/:id/like", middleware.JWTAuthMiddleware(), handler.ToggleLike)
}
