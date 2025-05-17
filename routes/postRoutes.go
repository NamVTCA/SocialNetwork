package routes

import (
    "socialnetwork/internal/post"
    "github.com/gin-gonic/gin"
    "socialnetwork/pkg/middleware"
)

func PostRoutes(r *gin.Engine, postHandler *post.PostHandler) {
    posts := r.Group("/posts")
    {
        posts.POST("", middleware.JWTAuthMiddleware(), postHandler.CreatePost)
        posts.GET("/:id", postHandler.GetPost)
        posts.PUT("/:id", middleware.JWTAuthMiddleware(), postHandler.UpdatePost)
        posts.DELETE("/:id", middleware.JWTAuthMiddleware(), postHandler.DeletePost)
        posts.GET("", postHandler.ListPosts)
    }
}
