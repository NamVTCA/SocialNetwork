package controllers

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func GetPosts(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "Danh sách bài viết",
    })
}
