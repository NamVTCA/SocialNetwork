package otp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"socialnetwork/dto/request"
	"socialnetwork/models"
)

type OTPHandler struct {
	otpService Service
}

func NewOTPHandler(otpService Service) *OTPHandler {
	return &OTPHandler{otpService: otpService}
}

func (h *OTPHandler) SendOTP(c *gin.Context) {
	var req request.SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.otpService.SendOTP(c.Request.Context(), &models.SendOTPRequest{
		Identifier: req.Identifier,
		Purpose:    req.Purpose,
		Channel:    req.Channel,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

func (h *OTPHandler) VerifyOTP(c *gin.Context) {
	var req request.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.otpService.VerifyOTP(c.Request.Context(), &models.VerifyOTPRequest{
		Identifier: req.Identifier,
		Purpose:    req.Purpose,
		OTP:        req.OTP,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
}
