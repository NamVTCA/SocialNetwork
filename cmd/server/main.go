package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"

	"socialnetwork/internal/otp"
	"socialnetwork/internal/user"
	"socialnetwork/pkg/config"
	"socialnetwork/pkg/email"
	"socialnetwork/pkg/sms"
	"socialnetwork/routes"
)

func main() {
	// Load biến môi trường từ .env
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
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisClient := redis.NewClient(&redis.Options{Addr: redisAddr})
	if _, err := redisClient.Ping(redisClient.Context()).Result(); err != nil {
		log.Fatalf("❌ Redis connection failed: %v", err)
	}

	// Khởi tạo email sender (SMTP hoặc mock)
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USERNAME")
	smtpPass := os.Getenv("SMTP_PASSWORD")

	var emailSender email.EmailSender
	if smtpHost != "" && smtpUser != "" && smtpPass != "" {
		emailSender = email.NewSMTPEmailSender(smtpHost, smtpPort, smtpUser, smtpPass)
		log.Println("✅ Using real SMTP email sender")
	} else {
		emailSender = email.NewMockEmailSender()
		log.Println("⚠️ SMTP config missing. Using mock email sender")
	}

	// Khởi tạo SMS sender (mock)
	smsSender := sms.NewMockSMSSender()
	log.Println("✅ Using MockSMS sender")

	// --- Khởi tạo services theo đúng thứ tự ---

	// 1. Repository người dùng
	userRepo := user.NewRepository(db)

	// 2. Service OTP (phải tạo trước userService vì userService cần dùng)
	otpService := otp.NewService(redisClient, emailSender, smsSender)

	// 3. Service User: truyền otpService vào đúng (thay vì nil)
	userService := user.NewService(userRepo, otpService, emailSender)

	// 4. Khởi tạo handler
	userHandler := user.NewHandler(userService)
	otpHandler := otp.NewOTPHandler(otpService)

	// --- Thiết lập Gin ---
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatalf("❌ Failed to set trusted proxies: %v", err)
	}

	// Đăng ký route
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	routes.UserRoutes(r, db, userHandler)
	routes.OTProutes(r, otpHandler)

	// Khởi động server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("🚀 Server is running at port:", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
