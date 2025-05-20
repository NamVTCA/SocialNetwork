package follow

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"socialnetwork/models"
)

// FollowRepository handles CRUD for Follow.
type FollowRepository interface {
	Create(ctx context.Context, follow *models.Follow) error
	Delete(ctx context.Context, follower, following primitive.ObjectID) error
	GetFollowers(ctx context.Context, userID primitive.ObjectID) ([]models.Follow, error)
	GetFollowing(ctx context.Context, userID primitive.ObjectID) ([]models.Follow, error)
	Exists(ctx context.Context, follower, following primitive.ObjectID) (bool, error)
}

// followRepo is concrete implementation.
type followRepo struct {
	col *mongo.Collection
}

// NewFollowRepository returns FollowRepository
func NewFollowRepository(db *mongo.Database) FollowRepository {
	return &followRepo{col: db.Collection("follows")}
}

func (r *followRepo) Create(ctx context.Context, follow *models.Follow) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	exists, err := r.Exists(ctx, follow.Follower, follow.Following)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("already following")
	}
	_, err = r.col.InsertOne(ctx, follow)
	return err
}

func (r *followRepo) Delete(ctx context.Context, follower, following primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.col.DeleteOne(ctx, bson.M{"follower": follower, "following": following})
	return err
}

func (r *followRepo) GetFollowers(ctx context.Context, userID primitive.ObjectID) ([]models.Follow, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.col.Find(ctx, bson.M{"following": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []models.Follow
	err = cursor.All(ctx, &list)
	return list, err
}

func (r *followRepo) GetFollowing(ctx context.Context, userID primitive.ObjectID) ([]models.Follow, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.col.Find(ctx, bson.M{"follower": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []models.Follow
	err = cursor.All(ctx, &list)
	return list, err
}

func (r *followRepo) Exists(ctx context.Context, follower, following primitive.ObjectID) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	count, err := r.col.CountDocuments(ctx, bson.M{"follower": follower, "following": following})
	return count > 0, err
}

