package follow

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"socialnetwork/models"
	
)

// FollowService defines follow logic.
type FollowService interface {
	Follow(ctx context.Context, follower, following primitive.ObjectID) error
	Unfollow(ctx context.Context, follower, following primitive.ObjectID) error
	GetFollowers(ctx context.Context, userID primitive.ObjectID) ([]models.Follow, error)
	GetFollowing(ctx context.Context, userID primitive.ObjectID) ([]models.Follow, error)
}

type followService struct {
	repo FollowRepository
}

// NewFollowService creates instance.
func NewFollowService(repo FollowRepository) FollowService {
	return &followService{repo: repo}
}

func (s *followService) Follow(ctx context.Context, follower, following primitive.ObjectID) error {
	f := &models.Follow{Follower: follower, Following: following, CreatedAt: time.Now()}
	return s.repo.Create(ctx, f)
}

func (s *followService) Unfollow(ctx context.Context, follower, following primitive.ObjectID) error {
	return s.repo.Delete(ctx, follower, following)
}

func (s *followService) GetFollowers(ctx context.Context, userID primitive.ObjectID) ([]models.Follow, error) {
	return s.repo.GetFollowers(ctx, userID)
}

func (s *followService) GetFollowing(ctx context.Context, userID primitive.ObjectID) ([]models.Follow, error) {
	return s.repo.GetFollowing(ctx, userID)
}