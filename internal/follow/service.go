package follow

import (
	"context"
	"errors"
	"socialnetwork/internal/notification"
	"socialnetwork/internal/user"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"socialnetwork/models"
)

type FollowService interface {
	FollowUser(ctx context.Context, followerID, followingID primitive.ObjectID) error
	UnfollowUser(ctx context.Context, followerID, followingID primitive.ObjectID) error
	GetFollowers(ctx context.Context, userID primitive.ObjectID) ([]primitive.ObjectID, error)
	GetFollowing(ctx context.Context, userID primitive.ObjectID) ([]primitive.ObjectID, error)
}

type followService struct {
	followRepo       FollowRepository
	userRepo         user.Repository
	notificationRepo notification.NotificationRepository
}

func NewFollowService(
	followRepo FollowRepository,
	userRepo user.Repository,
	notificationRepo notification.NotificationRepository,
) FollowService {
	return &followService{
		followRepo:       followRepo,
		userRepo:         userRepo,
		notificationRepo: notificationRepo,
	}
}

func (s *followService) FollowUser(ctx context.Context, followerID, followingID primitive.ObjectID) error {
	if followerID == followingID {
		return errors.New("cannot follow yourself")
	}

	isFollowing, err := s.followRepo.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return err
	}
	if isFollowing {
		return nil
	}

	err = s.followRepo.Create(ctx, &models.Follow{
		Follower:  followerID,
		Following: followingID,
	})
	if err != nil {
		return err
	}

	_ = s.userRepo.IncrementFollowerCount(ctx, followingID)
	_ = s.userRepo.IncrementFollowingCount(ctx, followerID)

	notif := &models.Notification{
		Recipient: followingID,
		Sender:    followerID,
		Type:      models.NotificationFollow,
		IsRead:    false,
		CreatedAt: time.Now(),
	}
	return s.notificationRepo.Create(ctx, notif)
}

func (s *followService) UnfollowUser(ctx context.Context, followerID, followingID primitive.ObjectID) error {
	err := s.followRepo.Delete(ctx, followerID, followingID)
	if err != nil {
		return err
	}
	_ = s.userRepo.DecrementFollowerCount(ctx, followingID)
	_ = s.userRepo.DecrementFollowingCount(ctx, followerID)
	return nil
}

func (s *followService) GetFollowers(ctx context.Context, userID primitive.ObjectID) ([]primitive.ObjectID, error) {
	followers, err := s.followRepo.GetFollowers(ctx, userID)
	if err != nil {
		return nil, err
	}

	followerIDs := make([]primitive.ObjectID, len(followers))
	for i, follow := range followers {
		followerIDs[i] = follow.Follower
	}
	return followerIDs, nil
}

func (s *followService) GetFollowing(ctx context.Context, userID primitive.ObjectID) ([]primitive.ObjectID, error) {
	following, err := s.followRepo.GetFollowing(ctx, userID)
	if err != nil {
		return nil, err
	}

	followingIDs := make([]primitive.ObjectID, len(following))
	for i, follow := range following {
		followingIDs[i] = follow.Following
	}
	return followingIDs, nil
}
