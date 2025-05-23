package post

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"socialnetwork/models"
	"strconv"
)

type PostHandler struct {
	postService PostService
}

func NewPostHandler(postService PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

// POST /posts
func (h *PostHandler) CreatePost(c *gin.Context) {
	var req models.Post
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// UserID nên lấy từ token auth (giả sử đã decode trước đó)
	// Ở đây mình giả định UserID đã được set trong context
	userIDStr := c.GetString("userID")
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id không hợp lệ"})
		return
	}
	req.UserID = userID

	post, err := h.postService.CreatePost(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, post)
}

// GET /posts/:id
func (h *PostHandler) GetPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id không hợp lệ"})
		return
	}

	post, err := h.postService.GetPostByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, post)
}

// PUT /posts/:id
func (h *PostHandler) UpdatePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id không hợp lệ"})
		return
	}

	// Lấy userID từ context (giả sử middleware đã lưu userID dạng string)
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Không tìm thấy thông tin người dùng"})
		return
	}
	userIDStr, ok := userIDValue.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Dữ liệu userID không hợp lệ"})
		return
	}
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UserID không hợp lệ"})
		return
	}

	// Lấy bài viết theo id để kiểm tra quyền
	post, err := h.postService.GetPostByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bài viết không tồn tại"})
		return
	}

	// Kiểm tra quyền sở hữu bài viết
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bạn không có quyền sửa bài viết này"})
		return
	}

	var req struct {
		Content  *string  `json:"content"`
		ImageURL *string  `json:"image_url"`
		Media    []string `json:"media"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateData := make(map[string]interface{})
	if req.Content != nil {
		updateData["content"] = *req.Content
	}
	if req.ImageURL != nil {
		updateData["image_url"] = *req.ImageURL
	}
	if req.Media != nil {
		updateData["media"] = req.Media
	}

	if len(updateData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "không có trường nào để cập nhật"})
		return
	}

	err = h.postService.UpdatePost(c.Request.Context(), id, updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cập nhật thành công"})
}

// DELETE /posts/:id
func (h *PostHandler) DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id không hợp lệ"})
		return
	}

	// Lấy userID từ token đã giải mã trong middleware
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Không xác thực được người dùng"})
		return
	}

	// Lấy bài viết từ DB
	post, err := h.postService.GetPostByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy bài viết"})
		return
	}

	// So sánh userID trong token và userID của bài viết
	if post.UserID.Hex() != userIDStr.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bạn không có quyền xóa bài viết này"})
		return
	}

	// Tiến hành xóa
	err = h.postService.DeletePost(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Xóa thành công"})
}

// GET /posts?page=1&limit=10
func (h *PostHandler) ListPosts(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit < 1 {
		limit = 10
	}

	posts, err := h.postService.ListPosts(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) GetPublicPostsByOwner(c *gin.Context) {
	ownerIDStr := c.Param("ownerID")
	ownerID, err := primitive.ObjectIDFromHex(ownerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ownerID"})
		return
	}

	posts, err := h.postService.GetPublicPostsByOwner(c.Request.Context(), ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get public posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}
