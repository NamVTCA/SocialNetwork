package short

import (
	"context"
	"time"

	"socialnetwork/internal/follow"
	"socialnetwork/internal/notification"
	"socialnetwork/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShortService interface {
	CreateShort(ctx context.Context, s *models.Short) error
	GetShortByID(ctx context.Context, id primitive.ObjectID) (*models.Short, error)
	GetShortsByOwner(ctx context.Context, ownerID primitive.ObjectID) ([]models.Short, error)
	IncrementView(ctx context.Context, id primitive.ObjectID) error
	DeleteShort(ctx context.Context, id, ownerID primitive.ObjectID) error
}

type shortService struct {
	repo             ShortRepository
	followRepo       follow.FollowRepository
	notificationRepo notification.NotificationRepository
}

func NewShortService(
	repo ShortRepository,
	followRepo follow.FollowRepository,
	notificationRepo notification.NotificationRepository,
) ShortService {
	return &shortService{repo, followRepo, notificationRepo}
}

func (s *shortService) CreateShort(ctx context.Context, sh *models.Short) error {
	err := s.repo.Create(ctx, sh)
	if err != nil {
		return err
	}

	followers, err := s.followRepo.GetFollowers(ctx, sh.OwnerID)
	if err == nil {
		for _, f := range followers {
			noti := &models.Notification{
				Recipient: f.Follower,
				Sender:    sh.OwnerID,
				Type:      models.NotificationNewShort,
				ShortID:   &sh.ID,
				IsRead:    false,
				CreatedAt: time.Now(),
			}
			_ = s.notificationRepo.Create(ctx, noti)
		}
	}

	return nil
}

func (s *shortService) GetShortByID(ctx context.Context, id primitive.ObjectID) (*models.Short, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *shortService) GetShortsByOwner(ctx context.Context, ownerID primitive.ObjectID) ([]models.Short, error) {
	return s.repo.GetByOwner(ctx, ownerID)
}

func (s *shortService) IncrementView(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.IncrementViews(ctx, id)
}

func (s *shortService) DeleteShort(ctx context.Context, id, ownerID primitive.ObjectID) error {
	return s.repo.Delete(ctx, id, ownerID)
}
