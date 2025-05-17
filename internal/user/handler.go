package user

import (
	"context"
	"net/http"
	"strconv"
	"fmt"
	"socialnetwork/dto/request"
	"socialnetwork/dto/response"
	"socialnetwork/models"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Register - Đăng ký tài khoản
func (h *Handler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 👈 log chi tiết lỗi
		return
	}

	user := &models.User{
		Username: req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.service.Register(context.TODO(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Đăng ký thành công"})
}

// Login - Đăng nhập
func (h *Handler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(context.TODO(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Đăng nhập thành công"})
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// GetMe - Lấy thông tin người dùng hiện tại
func (h *Handler) GetMe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.service.GetByID(context.TODO(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy thông tin người dùng"})
		return
	}

	resp := response.ToUserDetailResponse(user)
	c.JSON(http.StatusOK, resp)
}

// GetUsers - Lấy danh sách người dùng
func (h *Handler) GetUsers(c *gin.Context) {
	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}

	users, err := h.service.GetAllUsers(context.TODO())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lấy danh sách người dùng thất bại"})
		return
	}

	total := len(users)
	offset := (page - 1) * limit
	end := offset + limit
	if offset > total {
		offset = total
	}
	if end > total {
		end = total
	}

	var publicUsers []response.UserPublicResponse
	for _, u := range users[offset:end] {
		publicUsers = append(publicUsers, response.UserPublicResponse{
			ID:            u.ID.Hex(),
			Username:      u.Username,
			DisplayName:   u.DisplayName,
			AvatarURL:     u.AvatarURL,
			Bio:           u.Bio,
			FollowerCount: u.FollowerCount,
		})
	}

	resp := response.UserListResponse{
		Users: publicUsers,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	c.JSON(http.StatusOK, resp)
}


// UpateMe - Cập nhật thông tin cá nhân
func (h *Handler) UpdateMe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req request.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateProfile(c.Request.Context(), userID.(string), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật thông tin"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cập nhật thông tin thành công"})
}

// ChangePassword - Đổi mật khẩu
func (h *Handler) ChangePassword(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    var req request.ChangePasswordRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err := h.service.ChangePassword(c.Request.Context(), userID.(string), &req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Đổi mật khẩu thành công"})
}


func (h *Handler) ForgotPassword(c *gin.Context) {
	var req request.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.SendForgotPasswordOTP(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mã OTP đã được gửi vào email"})
}

func (h *Handler) ResetPassword(c *gin.Context) {
    var req request.ResetPasswordRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        fmt.Printf("BindJSON error: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    // fmt.Printf("Received ResetPasswordRequest: Email=%s, OTP=%s, NewPassword=****\n", req.Email, req.OTP)

    err := h.service.ResetPassword(c.Request.Context(), &req)
    if err != nil {
        fmt.Printf("ResetPassword error: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Đặt lại mật khẩu thành công"})
}

func (h *Handler) ChangeEmailRequest(c *gin.Context) {
	var req request.ChangeEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.ChangeEmailRequest(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Yêu cầu thay đổi email đã được gửi"})
}

func (h *Handler) VerifyEmailRequest(c *gin.Context) {
	var req request.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.VerifyEmailRequest(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Thay đổi email thành công"})
}