package short

import (
	"context"
	"time"

	"socialnetwork/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ShortRepository interface {
	Create(ctx context.Context, short *models.Short) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Short, error)
	GetByOwner(ctx context.Context, ownerID primitive.ObjectID) ([]models.Short, error)
	IncrementViews(ctx context.Context, id primitive.ObjectID) error
	Delete(ctx context.Context, id, ownerID primitive.ObjectID) error
	FindByOwnerAndVisibility(ctx context.Context, ownerID primitive.ObjectID, visibility string) ([]models.Short, error)
}

type shortRepository struct {
	collection *mongo.Collection
}

func NewShortRepository(db *mongo.Database) ShortRepository {
	return &shortRepository{collection: db.Collection("shorts")}
}

func (r *shortRepository) Create(ctx context.Context, short *models.Short) error {
	short.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, short)
	return err
}

func (r *shortRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Short, error) {
	var s models.Short
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&s)
	return &s, err
}

func (r *shortRepository) GetByOwner(ctx context.Context, ownerID primitive.ObjectID) ([]models.Short, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"ownerId": ownerID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var shorts []models.Short
	err = cursor.All(ctx, &shorts)
	return shorts, err
}

func (r *shortRepository) IncrementViews(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"views": 1}})
	return err
}

func (r *shortRepository) Delete(ctx context.Context, id, ownerID primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id, "ownerId": ownerID})
	return err
}

func (r *shortRepository) FindByOwnerAndVisibility(ctx context.Context, ownerID primitive.ObjectID, visibility string) ([]models.Short, error) {
    filter := bson.M{
        "owner_id":   ownerID,
        "visibility": visibility,
    }
    cursor, err := r.collection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }
    var shorts []models.Short
    if err = cursor.All(ctx, &shorts); err != nil {
        return nil, err
    }
    return shorts, nil
}
