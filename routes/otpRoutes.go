package routes

import (
	"github.com/gin-gonic/gin"
	"socialnetwork/internal/otp"
)

func OTProutes(r *gin.Engine, otpHandler *otp.OTPHandler) {
	otpGroup := r.Group("/otp")
	{
		otpGroup.POST("/sendEmail", otpHandler.SendOTP)
		otpGroup.POST("/verifyEmail", otpHandler.VerifyOTP)
		otpGroup.POST("/sendSMS", otpHandler.SendOTP)
		otpGroup.POST("/verifySMS", otpHandler.VerifyOTP)

	}
}
