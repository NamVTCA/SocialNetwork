package comment

import (
    "net/http"

    "socialnetwork/models"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentHandler struct {
    service Service
}

func NewCommentHandler(service Service) *CommentHandler {
    return &CommentHandler{service: service}
}

// POST /comments/:postID
func (h *CommentHandler) CreateComment(c *gin.Context) {
    userIDStr, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    postIDStr := c.Param("postID")
    postID, err := primitive.ObjectIDFromHex(postIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
        return
    }

    var req struct {
        Content string `json:"content" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    comment := &models.Comment{
        PostID:  postID,
        UserID:  userID,
        Content: req.Content,
    }

    created, err := h.service.Create(c.Request.Context(), comment)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, created)
}

// GET /comments/:postID
func (h *CommentHandler) GetCommentsByPost(c *gin.Context) {
    postIDStr := c.Param("postID")
    postID, err := primitive.ObjectIDFromHex(postIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
        return
    }

    comments, err := h.service.ListByPost(c.Request.Context(), postID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, comments)
}

// PUT /comments/:id
func (h *CommentHandler) UpdateComment(c *gin.Context) {
    commentIDStr := c.Param("id")
    commentID, err := primitive.ObjectIDFromHex(commentIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
        return
    }

    userIDStr, _ := c.Get("userID")
    userID, _ := primitive.ObjectIDFromHex(userIDStr.(string))

    var req models.CommentUpdateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    update := make(map[string]interface{})
    if req.Content != nil {
        update["content"] = *req.Content
    }

    if len(update) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Nothing to update"})
        return
    }

    err = h.service.Update(c.Request.Context(), commentID, userID, update)
    if err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Comment updated"})
}

// DELETE /comments/:id
func (h *CommentHandler) DeleteComment(c *gin.Context) {
    commentIDStr := c.Param("id")
    commentID, err := primitive.ObjectIDFromHex(commentIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
        return
    }

    userIDStr, _ := c.Get("userID")
    userID, _ := primitive.ObjectIDFromHex(userIDStr.(string))

    err = h.service.Delete(c.Request.Context(), commentID, userID)
    if err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}

// PUT /comments/:id/like
func (h *CommentHandler) ToggleLike(c *gin.Context) {
    commentIDStr := c.Param("id")
    commentID, err := primitive.ObjectIDFromHex(commentIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
        return
    }

    userIDStr, _ := c.Get("userID")
    userID, _ := primitive.ObjectIDFromHex(userIDStr.(string))

    err = h.service.ToggleLike(c.Request.Context(), commentID, userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Like status updated"})
}
