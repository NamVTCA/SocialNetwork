package models

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
    ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    UserID    primitive.ObjectID   `bson:"user_id" json:"user_id"`
    Content   string               `bson:"content" json:"content"`
    ImageURL  string               `bson:"image_url,omitempty" json:"image_url,omitempty"`
    Media     []string             `bson:"media,omitempty" json:"media,omitempty"` // ảnh/video (optional)
    Likes     []primitive.ObjectID `bson:"likes,omitempty" json:"likes,omitempty"` // danh sách user đã like
    Dislikes  []primitive.ObjectID `bson:"dislikes,omitempty" json:"dislikes,omitempty"` // danh sách user đã dislike
    Visibility string               `bson:"visibility" json:"visibility"`           // "public", "private", "friends"...
    CreatedAt time.Time            `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time            `bson:"updated_at" json:"updated_at"`
}

type PostUpdateRequest struct {
    Content  *string  `json:"content"`
    ImageURL *string  `json:"image_url"`
    Media    []string `json:"media"`
}