package auth

import (
    "fmt"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

// Tạo JWT từ user ID
func GenerateJWT(userID string) (string, error) {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        return "", fmt.Errorf("JWT_SECRET không được để trống")
    }

    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
        "iat":     time.Now().Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

// Giải mã và xác thực token
func ValidateJWT(tokenStr string) (string, error) {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        return "", fmt.Errorf("JWT_SECRET không được để trống")
    }

    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        // Kiểm tra thuật toán ký
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("phương thức ký không hợp lệ: %v", token.Header["alg"])
        }
        return []byte(secret), nil
    })

    if err != nil || !token.Valid {
        return "", fmt.Errorf("token không hợp lệ: %v", err)
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        return "", fmt.Errorf("claims không hợp lệ")
    }

    userID, ok := claims["user_id"].(string)
    if !ok {
        return "", fmt.Errorf("không tìm thấy user_id trong token")
    }

    return userID, nil
}
