package user

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"socialnetwork/dto/request"
	"socialnetwork/dto/response"
	"socialnetwork/models"
	"strconv"
	"strings"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func normalizePhone(phone string) string {
	if strings.HasPrefix(phone, "0") {
		return "+84" + phone[1:]
	}
	return phone
}

// Register - Đăng ký người dùng mới
func (h *Handler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Khởi tạo user cơ bản
	user := &models.User{
		Username: req.Name,
		Password: req.Password, // sẽ hash ở service
	}

	// Xác định xem identifier là email hay số điện thoại
	if strings.Contains(req.Identifier, "@") {
		user.Email = req.Identifier
	} else {
		user.Phone = normalizePhone(req.Identifier)
	}

	// Gọi service
	if err := h.service.Register(context.TODO(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Đăng ký thành công",
	})
}

// Login - Đăng nhập
func (h *Handler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(context.TODO(), normalizePhone(req.Identifier), req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Đăng nhập thành công",
		"token":   token,
	})
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
		// fmt.Printf("ResetPassword error: %v\n", err)
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

	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, ok := userIDRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	err := h.service.ChangeEmailRequest(c.Request.Context(), userID, &req)
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

	// ✅ Lấy userID từ context do middleware JWT gán
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, ok := userIDRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	err := h.service.VerifyEmailRequest(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Thay đổi email thành công"})
}

// POST /users/:id/friends/request
func (h *Handler) SendFriendRequest(c *gin.Context) {
	userIDStr := c.GetString("userID")
	fromID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sender ID"})
		return
	}

	toIDStr := c.Param("id")
	toID, err := primitive.ObjectIDFromHex(toIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid receiver ID"})
		return
	}

	// 1. Không cho tự gửi kết bạn cho chính mình
	if fromID == toID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không thể gửi kết bạn cho chính mình"})
		return
	}

	// 2. Kiểm tra nếu đã gửi yêu cầu rồi
	exists, err := h.service.FriendRequestExists(c.Request.Context(), fromID, toID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không kiểm tra được trạng thái lời mời"})
		return
	}

	if exists {
		// 3. Nếu đã gửi rồi thì thực hiện hủy lời mời
		err := h.service.CancelFriendRequest(c.Request.Context(), fromID, toID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Hủy lời mời thất bại"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Đã hủy lời mời kết bạn"})
		return
	}

	// Gửi yêu cầu kết bạn mới
	if err := h.service.SendFriendRequest(c.Request.Context(), fromID, toID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Yêu cầu kết bạn đã được gửi"})
}

// POST /users/:id/friends/accept
func (h *Handler) AcceptFriendRequest(c *gin.Context) {
	// Người đang đăng nhập (phải là người được gửi lời mời)
	receiverIDStr := c.GetString("userID")
	if receiverIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	receiverID, err := primitive.ObjectIDFromHex(receiverIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID in token"})
		return
	}

	// Người đã gửi lời mời (truyền qua URL)
	senderIDStr := c.Param("id")
	senderID, err := primitive.ObjectIDFromHex(senderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID in URL"})
		return
	}

	// Gọi service xử lý
	if err := h.service.AcceptFriendRequest(c.Request.Context(), receiverID, senderID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đã chấp nhận yêu cầu kết bạn"})
}

// POST /users/:id/block
func (h *Handler) BlockUser(c *gin.Context) {
	userIDStr := c.GetString("userID")
	uid, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	targetIDStr := c.Param("id")
	targetID, err := primitive.ObjectIDFromHex(targetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID"})
		return
	}

	if err := h.service.BlockUser(c.Request.Context(), uid, targetID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đã chặn người dùng"})
}

// PUT /users/me/hide-profile
func (h *Handler) ToggleHideProfile(c *gin.Context) {
	userID := c.GetString("userID")
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var body struct {
		Hide bool `json:"hide"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ToggleHideProfile(c.Request.Context(), uid, body.Hide); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	status := "hiện"
	if body.Hide {
		status = "ẩn"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Đã cập nhật chế độ trang cá nhân thành: %s", status),
		"hide":    body.Hide, // thêm trạng thái trả về nếu cần
	})
}
