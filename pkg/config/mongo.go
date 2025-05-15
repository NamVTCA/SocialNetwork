package config

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var Database *mongo.Database

func ConnectMongoDB() (*mongo.Database, error) {
    // Load biến môi trường từ file .env
    err := godotenv.Load()
    if err != nil {
        log.Println("⚠️ Không thể load file .env, dùng biến môi trường hiện tại")
    }

    uri := os.Getenv("MONGO_URI")
    dbName := os.Getenv("MONGO_DB_NAME")

    if uri == "" || dbName == "" {
        return nil, fmt.Errorf("❌ MONGO_URI hoặc MONGO_DB_NAME không được để trống")
    }

    clientOpts := options.Client().ApplyURI(uri)

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, clientOpts)
    if err != nil {
        return nil, fmt.Errorf("❌ Kết nối MongoDB thất bại: %v", err)
    }

    // Ping kiểm tra kết nối
    err = client.Ping(ctx, nil)
    if err != nil {
        return nil, fmt.Errorf("❌ Không thể ping MongoDB: %v", err)
    }

    log.Println("✅ Đã kết nối MongoDB!")

    Database = client.Database(dbName)
    return Database, nil
}
