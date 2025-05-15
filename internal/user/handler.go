package user

import (
	"context"
	"net/http"
	"strconv"

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

	c.JSON(http.StatusOK, gin.H{"token": token})
}

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

func (h *Handler) GetUsers(c *gin.Context) {
	// Mặc định phân trang
	page := 1
	limit := 10

	// Đọc tham số truy vấn ?page=...&limit=...
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

	offset := (page - 1) * limit

	// Lấy tất cả người dùng
	users, err := h.service.GetAllUsers()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Tổng số user
	total := len(users)

	// Lọc ra theo page/limit
	end := offset + limit
	if end > total {
		end = total
	}
	if offset > total {
		offset = total
	}

	usersPage := users[offset:end]

	// Chuyển sang DTO
	var userResponses []response.UserPublicResponse
	for _, u := range usersPage {
		userResponses = append(userResponses, response.UserPublicResponse{
			ID:            u.ID.Hex(),
			Username:      u.Username,
			DisplayName:   u.DisplayName,
			AvatarURL:     u.AvatarURL,
			Bio:           u.Bio,
			FollowerCount: u.FollowerCount,
		})
	}


	c.JSON(http.StatusOK, response.UserListResponse{
		Users: userResponses,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}
