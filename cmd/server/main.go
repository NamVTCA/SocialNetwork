package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"

	"socialnetwork/internal/comment"
	"socialnetwork/internal/follow"
	"socialnetwork/internal/notification"
	"socialnetwork/internal/otp"
	"socialnetwork/internal/post"
	"socialnetwork/internal/short"
	"socialnetwork/internal/user"
	"socialnetwork/internal/video"
	"socialnetwork/pkg/config"
	"socialnetwork/pkg/email"
	"socialnetwork/pkg/sms"
	"socialnetwork/routes"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Không tìm thấy file .env hoặc lỗi khi load, dùng biến môi trường hệ thống")
	}

	// MongoDB
	db, err := config.ConnectMongoDB()
	if err != nil {
		log.Fatalf("❌ MongoDB connection failed: %v", err)
	}

	// Redis
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisClient := redis.NewClient(&redis.Options{Addr: redisAddr})
	if _, err := redisClient.Ping(redisClient.Context()).Result(); err != nil {
		log.Fatalf("❌ Redis connection failed: %v", err)
	}

	// Email & SMS Sender
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

	smsSender := sms.NewMockSMSSender()
	log.Println("✅ Using MockSMS sender")

	// Services & Handlers
	userRepo := user.NewRepository(db)
	otpService := otp.NewService(redisClient, emailSender, smsSender, userRepo)
	userService := user.NewService(userRepo, otpService, emailSender)

	userHandler := user.NewHandler(userService)
	otpHandler := otp.NewOTPHandler(otpService)

	postRepo := post.NewPostRepository(db)
	postService := post.NewPostService(&postRepo)
	postHandler := post.NewPostHandler(postService)

	commentRepo := comment.NewCommentRepository(db)
	commentService := comment.NewCommentService(commentRepo)
	commentHandler := comment.NewCommentHandler(commentService)

	notifRepo := notification.NewNotificationRepository(db)
	notifService := notification.NewNotificationService(notifRepo)
	notifHandler := notification.NewNotificationHandler(notifService)

	followRepo := follow.NewFollowRepository(db)
	followService := follow.NewFollowService(followRepo, userRepo, notifRepo)
	followHandler := follow.NewFollowHandler(followService)

	videoRepo := video.NewVideoRepository(db)
	videoService := video.NewVideoService(videoRepo, followRepo, notifRepo)

	shortRepo := short.NewShortRepository(db)
	shortService := short.NewShortService(shortRepo, followRepo, notifRepo)

	// Gin Setup
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// ✅ CORS Middleware cho phép React kết nối
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // FE Vite default
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatalf("❌ Failed to set trusted proxies: %v", err)
	}

	// Routes
	r.POST("/registerEmail", userHandler.RegisterEmail)
	r.POST("/registerPhone", userHandler.RegisterPhone)
	r.POST("/loginEmail", userHandler.LoginEmail)
	r.POST("/loginPhone", userHandler.LoginPhone)

	routes.UserRoutes(r, db, userHandler)
	routes.OTProutes(r, otpHandler)
	routes.PostRoutes(r, postHandler)
	routes.CommentRoutes(r, db, commentHandler)

	api := r.Group("/api")
	routes.FollowRoutes(api, db, followHandler)
	routes.Video_ShortRoutes(r, videoService, shortService)
	routes.NotificationRoutes(api, notifHandler)

	// Run server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("🚀 Server is running at port:", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
