package follow

import (
	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"socialnetwork/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FollowRepository interface {
	Create(ctx context.Context, follow *models.Follow) error
	Delete(ctx context.Context, followerID, followingID primitive.ObjectID) error
	IsFollowing(ctx context.Context, followerID, followingID primitive.ObjectID) (bool, error)
	GetFollowers(ctx context.Context, userID primitive.ObjectID) ([]models.Follow, error)
	GetFollowing(ctx context.Context, userID primitive.ObjectID) ([]models.Follow, error)
}

type followRepository struct {
	collection *mongo.Collection
}

func NewFollowRepository(db *mongo.Database) FollowRepository {
	return &followRepository{
		collection: db.Collection("follows"),
	}
}

func (r *followRepository) Create(ctx context.Context, follow *models.Follow) error {
	follow.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, follow)
	return err
}

func (r *followRepository) Delete(ctx context.Context, followerID, followingID primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{
		"follower":  followerID,
		"following": followingID,
	})
	return err
}

func (r *followRepository) IsFollowing(ctx context.Context, followerID, followingID primitive.ObjectID) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"follower":  followerID,
		"following": followingID,
	})
	return count > 0, err
}

func (r *followRepository) GetFollowers(ctx context.Context, userID primitive.ObjectID) ([]models.Follow, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"following": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var follows []models.Follow
	if err := cursor.All(ctx, &follows); err != nil {
		return nil, err
	}
	return follows, nil
}

func (r *followRepository) GetFollowing(ctx context.Context, userID primitive.ObjectID) ([]models.Follow, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"follower": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var follows []models.Follow
	if err := cursor.All(ctx, &follows); err != nil {
		return nil, err
	}
	return follows, nil
}
