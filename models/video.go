package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Video struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	OwnerID     primitive.ObjectID `bson:"ownerId" json:"ownerId"` // user đăng
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	URL         string             `bson:"url" json:"url"`
	Thumbnail   string             `bson:"thumbnail,omitempty" json:"thumbnail,omitempty"`
	Duration    int                `bson:"duration,omitempty" json:"duration,omitempty"` // giây
	Likes       int                `bson:"likes" json:"likes"`
	Views       int                `bson:"views" json:"views"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}
