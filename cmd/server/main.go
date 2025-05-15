package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv" 

	"socialnetwork/internal/user"
	"socialnetwork/pkg/config"
	"socialnetwork/routes"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env, using system env vars")
	}

	gin.SetMode(gin.ReleaseMode)

	db, err := config.ConnectMongoDB()
	if err != nil {
		log.Fatalf("MongoDB connection failed:❌ %v", err)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatalf("❌ Failed to set trusted proxies: %v", err)
	}

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	// Đăng ký route đăng ký / đăng nhập
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	// Đăng ký các route khác như /users/, /users/me
	routes.UserRoutes(r, db, userHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("🚀 Server is running at port:", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
