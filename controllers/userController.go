package controllers

import (
	"net/http"


	"github.com/gin-gonic/gin"
)



func GetUsers(c *gin.Context) {
	// Fake demo response
	c.JSON(http.StatusOK, gin.H{
		"message": "Danh sách người dùng",

	})
}

func GetMe(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "User information retrieved successfully",
    })
}


