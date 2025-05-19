package notification


import (
	"context"
	"time"
	"socialnetwork/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationService interface {
	GetUserNotifications(ctx context.Context, userID primitive.ObjectID) ([]models.Notification, error)
	CreateNotification(ctx context.Context, notif *models.Notification) error
	ReadNotification(ctx context.Context, notifID primitive.ObjectID) error
}

type notificationService struct {
	repo NotificationRepository
}

func NewNotificationService(repo NotificationRepository) NotificationService {
	return &notificationService{repo: repo}
}

func (s *notificationService) GetUserNotifications(ctx context.Context, userID primitive.ObjectID) ([]models.Notification, error) {
	return s.repo.GetByRecipient(ctx, userID)
}

func (s *notificationService) CreateNotification(ctx context.Context, notif *models.Notification) error {
	return s.repo.Create(ctx, notif)
}

func (s *notificationService) ReadNotification(ctx context.Context, notifID primitive.ObjectID) error {
	return s.repo.MarkAsRead(ctx, notifID)
}



func (s *notificationService) NotifyComment(ctx context.Context, senderID, postOwnerID, postID primitive.ObjectID, message string) error {
	if senderID == postOwnerID {
		return nil // không tự thông báo cho mình
	}
	notif := &models.Notification{
		ID:        primitive.NewObjectID(),
		Recipient: postOwnerID,
		Sender:    senderID,
		Type:      models.NotificationComment,
		PostID:    &postID,
		Message:   message,
		IsRead:    false,
		CreatedAt: time.Now(),
	}
	return s.repo.Create(ctx, notif)
}

func (s *notificationService) NotifyNewContent(ctx context.Context, senderID primitive.ObjectID, followerIDs []primitive.ObjectID, postID primitive.ObjectID, contentType models.NotificationType) error {
	for _, followerID := range followerIDs {
		if followerID == senderID {
			continue
		}
		notif := &models.Notification{
			ID:        primitive.NewObjectID(),
			Recipient: followerID,
			Sender:    senderID,
			Type:      contentType,
			PostID:    &postID,
			IsRead:    false,
			CreatedAt: time.Now(),
		}
		if err := s.repo.Create(ctx, notif); err != nil {
			return err
		}
	}
	return nil
}