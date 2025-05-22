package video

import (
	"context"
	"time"

	"socialnetwork/internal/notification"
	"socialnetwork/internal/follow"
	"socialnetwork/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VideoService interface {
	CreateVideo(ctx context.Context, video *models.Video) error
	GetVideoByID(ctx context.Context, id primitive.ObjectID) (*models.Video, error)
	GetVideosByOwner(ctx context.Context, ownerID primitive.ObjectID) ([]models.Video, error)
	IncrementView(ctx context.Context, id primitive.ObjectID) error
	DeleteVideo(ctx context.Context, id, ownerID primitive.ObjectID) error
}

type videoService struct {
	videoRepo        VideoRepository
	followRepo       follow.FollowRepository
	notificationRepo notification.NotificationRepository
}

func NewVideoService(
	videoRepo VideoRepository,
	followRepo follow.FollowRepository,
	notificationRepo notification.NotificationRepository,
) VideoService {
	return &videoService{videoRepo, followRepo, notificationRepo}
}

func (s *videoService) CreateVideo(ctx context.Context, video *models.Video) error {
	err := s.videoRepo.Create(ctx, video)
	if err != nil {
		return err
	}

	// Gửi thông báo đến followers
	followers, err := s.followRepo.GetFollowers(ctx, video.OwnerID)
	if err == nil {
		for _, f := range followers {
			noti := &models.Notification{
				Recipient: f.Follower,
				Sender:    video.OwnerID,
				Type:      models.NotificationNewVideo,
				VideoID:   &video.ID,
				IsRead:    false,
				CreatedAt: time.Now(),
			}
			_ = s.notificationRepo.Create(ctx, noti)
		}
	}

	return nil
}

func (s *videoService) GetVideoByID(ctx context.Context, id primitive.ObjectID) (*models.Video, error) {
	return s.videoRepo.GetByID(ctx, id)
}

func (s *videoService) GetVideosByOwner(ctx context.Context, ownerID primitive.ObjectID) ([]models.Video, error) {
	return s.videoRepo.GetByOwner(ctx, ownerID)
}

func (s *videoService) IncrementView(ctx context.Context, id primitive.ObjectID) error {
	return s.videoRepo.IncrementViews(ctx, id)
}

func (s *videoService) DeleteVideo(ctx context.Context, id, ownerID primitive.ObjectID) error {
	return s.videoRepo.Delete(ctx, id, ownerID)
}
