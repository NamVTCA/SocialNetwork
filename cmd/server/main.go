package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"

	"socialnetwork/internal/user"
	"socialnetwork/internal/otp"
	"socialnetwork/pkg/config"
	"socialnetwork/pkg/email"
	"socialnetwork/pkg/sms"
	"socialnetwork/routes"
)

func main() {
	// Load biến môi trường
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ Không tìm thấy file .env hoặc lỗi khi load, dùng biến môi trường hệ thống")
	}

	// Kết nối MongoDB
	db, err := config.ConnectMongoDB()
	if err != nil {
		log.Fatalf("❌ MongoDB connection failed: %v", err)
	}

	// Kết nối Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"), // ví dụ "localhost:6379"
	})
	_, err = redisClient.Ping(redisClient.Context()).Result()
	if err != nil {
		log.Fatalf("❌ Redis connection failed: %v", err)
	}

	// Email Sender
	emailSender := email.NewMockEmailSender() // hoặc dùng SMTP thật

	// SMS Sender
	smsSender := sms.NewMockSMSSender() // hoặc tích hợp Twilio thật

	// Init Service & Handler cho User
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	// Init Service & Handler cho OTP
	otpService := otp.NewService(redisClient, emailSender, smsSender)
	otpHandler := otp.NewOTPHandler(otpService)

	// Init router
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatalf("❌ Failed to set trusted proxies: %v", err)
	}

	// Đăng ký route người dùng
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	routes.UserRoutes(r, db, userHandler)

	// Đăng ký route OTP
	routes.OTProutes(r, otpHandler)

	// Cổng server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("🚀 Server is running at port:", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
