package routes

import (
	"github.com/gin-gonic/gin"
	"socialnetwork/pkg/middleware"
	"socialnetwork/internal/short"
	"socialnetwork/internal/video"
)

func Video_ShortRoutes(r *gin.Engine, videoService video.VideoService, shortService short.ShortService) {
	videoHandler := video.NewVideoHandler(videoService)
	shortHandler := short.NewShortHandler(shortService)

	authRoutes := r.Group("/api", middleware.JWTAuthMiddleware())
	{
		// Video routes
		authRoutes.POST("/videos", videoHandler.CreateVideo)
		authRoutes.GET("/videos/:id", videoHandler.GetVideoByID)
		authRoutes.GET("/videos", videoHandler.GetVideosByOwner)
		authRoutes.DELETE("/videos/:id", videoHandler.DeleteVideo)

		// Short routes
		authRoutes.POST("/shorts", shortHandler.CreateShort)
		authRoutes.GET("/shorts/:id", shortHandler.GetShortByID)
		authRoutes.GET("/shorts", shortHandler.GetShortsByOwner)
		authRoutes.DELETE("/shorts/:id", shortHandler.DeleteShort)
	}
}
