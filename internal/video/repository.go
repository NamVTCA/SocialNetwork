package video

import (
	"context"
	"time"

	"socialnetwork/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VideoRepository interface {
	Create(ctx context.Context, video *models.Video) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Video, error)
	GetByOwner(ctx context.Context, ownerID primitive.ObjectID) ([]models.Video, error)
	IncrementViews(ctx context.Context, id primitive.ObjectID) error
	Delete(ctx context.Context, id primitive.ObjectID, ownerID primitive.ObjectID) error
	FindByOwnerAndVisibility(ctx context.Context, ownerID primitive.ObjectID, visibility string) ([]models.Video, error)
}

type videoRepository struct {
	collection *mongo.Collection
}

func NewVideoRepository(db *mongo.Database) VideoRepository {
	return &videoRepository{collection: db.Collection("videos")}
}

func (r *videoRepository) Create(ctx context.Context, video *models.Video) error {
	video.CreatedAt = time.Now()
	video.UpdatedAt = video.CreatedAt
	_, err := r.collection.InsertOne(ctx, video)
	return err
}

func (r *videoRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Video, error) {
	var video models.Video
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&video)
	return &video, err
}

func (r *videoRepository) GetByOwner(ctx context.Context, ownerID primitive.ObjectID) ([]models.Video, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"ownerId": ownerID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var videos []models.Video
	err = cursor.All(ctx, &videos)
	return videos, err
}

func (r *videoRepository) IncrementViews(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"views": 1}})
	return err
}

func (r *videoRepository) Delete(ctx context.Context, id, ownerID primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id, "ownerId": ownerID})
	return err
}

func (r *videoRepository) FindByOwnerAndVisibility(ctx context.Context, ownerID primitive.ObjectID, visibility string) ([]models.Video, error) {
    filter := bson.M{
        "owner_id":  ownerID,
        "visibility": visibility,
    }
    cursor, err := r.collection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }
    var videos []models.Video
    if err = cursor.All(ctx, &videos); err != nil {
        return nil, err
    }
    return videos, nil
}
