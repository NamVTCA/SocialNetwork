package routes

import (
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/mongo"
    "socialnetwork/controllers"
)

func PostRoutes(r *gin.Engine, db *mongo.Client) {
    postRoutes := r.Group("/posts")
    {
        postRoutes.GET("/", controllers.GetPosts)
    }
}
