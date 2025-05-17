package models

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
    ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    PostID    primitive.ObjectID   `bson:"post_id" json:"post_id"`
    UserID    primitive.ObjectID   `bson:"user_id" json:"user_id"`
    Content   string               `bson:"content" json:"content"`
    Likes     []primitive.ObjectID `bson:"likes,omitempty" json:"likes,omitempty"`
    CreatedAt time.Time            `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time            `bson:"updated_at" json:"updated_at"`
}

type CommentUpdateRequest struct {
    Content *string `json:"content"`
	Likes   []primitive.ObjectID `json:"likes,omitempty"`
}