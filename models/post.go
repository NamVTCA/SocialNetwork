package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
    ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    UserID    primitive.ObjectID   `bson:"user_id" json:"user_id"`
    Content   string               `bson:"content" json:"content"`
    ImageURL  string               `bson:"image_url,omitempty" json:"image_url,omitempty"`
    Likes     []primitive.ObjectID `bson:"likes" json:"likes"`
    CreatedAt int64                `bson:"created_at" json:"created_at"`
}
