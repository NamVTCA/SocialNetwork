package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type NotificationType string

const (
	NotificationLike     NotificationType = "like"
	NotificationComment  NotificationType = "comment"     // Ai đó comment bài viết của bạn
	NotificationFollow   NotificationType = "follow"      // Ai đó theo dõi bạn
	NotificationMention  NotificationType = "mention"     // Bạn được nhắc tên
	NotificationSystem   NotificationType = "system"      // Thông báo hệ thống
	NotificationNewPost  NotificationType = "new_post"    // Người bạn theo dõi tạo post
	NotificationNewVideo NotificationType = "new_video"   // Người bạn theo dõi tạo video
	NotificationNewShort NotificationType = "new_short"   // Người bạn theo dõi tạo short
)

type Notification struct {
	ID         primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	Recipient  primitive.ObjectID  `bson:"recipient" json:"recipient"`       // Người nhận
	Sender     primitive.ObjectID  `bson:"sender,omitempty" json:"sender,omitempty"` // Người tạo hành động (người comment, người tạo bài viết,...)
	Type       NotificationType    `bson:"type" json:"type"`
	PostID     *primitive.ObjectID `bson:"post,omitempty" json:"post,omitempty"`     // Bài viết liên quan
	VideoID    *primitive.ObjectID `bson:"video,omitempty" json:"video,omitempty"`   // Video liên quan
	ShortID    *primitive.ObjectID `bson:"short,omitempty" json:"short,omitempty"`   // Short liên quan
	Message    string              `bson:"message,omitempty" json:"message,omitempty"` // Thông báo tuỳ chỉnh
	IsRead     bool                `bson:"isRead" json:"isRead"`
	CreatedAt  time.Time           `bson:"createdAt" json:"createdAt"`
}
