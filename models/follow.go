package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Follow struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Follower  primitive.ObjectID `bson:"follower" json:"follower"`   // Người theo dõi (người bấm follow)
	Following primitive.ObjectID `bson:"following" json:"following"` // Người được theo dõi
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}
