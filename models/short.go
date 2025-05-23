package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Short struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	OwnerID   primitive.ObjectID `bson:"ownerId" json:"ownerId"`
	Title     string             `bson:"title" json:"title"`
	URL       string             `bson:"url" json:"url"`
	Duration  int                `bson:"duration" json:"duration"` // thường < 60s
	Likes     int                `bson:"likes" json:"likes"`
	Dislikes  int                `bson:"dislikes" json:"dislikes"`
	Visibility string               `bson:"visibility" json:"visibility"`
	Views     int                `bson:"views" json:"views"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}
