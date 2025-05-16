package routes

import (
	"github.com/gin-gonic/gin"
	"socialnetwork/internal/otp"
)

func OTProutes(r *gin.Engine, otpHandler *otp.OTPHandler) {
	otpGroup := r.Group("/otp")
	{
		otpGroup.POST("/send", otpHandler.SendOTP)
		otpGroup.POST("/verify", otpHandler.VerifyOTP)
		otpGroup.GET("/debugget", otpHandler.DebugGetOTP)
	}
}
