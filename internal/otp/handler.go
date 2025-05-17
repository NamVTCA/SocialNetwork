package otp

import (
	"net/http"
	// "fmt"
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
		// fmt.Printf("[DEBUG][SendOTP] Invalid request body: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// fmt.Printf("[DEBUG][SendOTP] Received request: identifier=%s, purpose=%s, channel=%s\n", req.Identifier, req.Purpose, req.Channel)

	err := h.otpService.SendOTP(c.Request.Context(), &models.SendOTPRequest{
		Identifier: req.Identifier,
		Purpose:    req.Purpose,
		Channel:    req.Channel,
	})
	if err != nil {
		// fmt.Printf("[DEBUG][SendOTP] Failed to send OTP: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// fmt.Printf("[DEBUG][SendOTP] OTP sent successfully for identifier=%s\n", req.Identifier)
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

func (h *OTPHandler) VerifyOTP(c *gin.Context) {
	var req request.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// fmt.Printf("[DEBUG][VerifyOTP] Invalid request body: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// fmt.Printf("[DEBUG][VerifyOTP] Received request: identifier=%s, purpose=%s, otp=%s\n", req.Identifier, req.Purpose, req.OTP)

	err := h.otpService.VerifyOTP(c.Request.Context(), &models.VerifyOTPRequest{
		Identifier: req.Identifier,
		Purpose:    req.Purpose,
		OTP:        req.OTP,
	})
	if err != nil {
		// fmt.Printf("[DEBUG][VerifyOTP] OTP verification failed: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// fmt.Printf("[DEBUG][VerifyOTP] OTP verified successfully for identifier=%s\n", req.Identifier)
	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
}
