package otp

import (
	"net/http"
	"fmt"
	"github.com/gin-gonic/gin"
)
func (h *OTPHandler) DebugGetOTP(c *gin.Context) {
	type request struct {
		Identifier string `json:"identifier" binding:"required,email"`
		Purpose    string `json:"purpose" binding:"required"`
	}

	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	key := fmt.Sprintf("otp:%s:%s", req.Identifier, req.Purpose)
	otp, err := h.otpService.(*service).redisClient.Get(c.Request.Context(), key).Result()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "OTP not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"otp": otp})
}
